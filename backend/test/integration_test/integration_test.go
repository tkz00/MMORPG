package integration_test

import (
	"backend/pkg/server"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/gorilla/websocket"
)

// ----------------------
// Global setup/teardown
// ----------------------

func TestMain(m *testing.M) {
	os.Setenv("GO_ENV", "test")

	// Try to load .env only if it exists
	cwd, _ := os.Getwd()
	root := filepath.Join(cwd, "..", "..", "..")
	envPath := filepath.Join(root, ".env")

	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			fmt.Println("Warning: failed to load .env:", err)
		} else {
			fmt.Println("Loaded local .env for tests")
		}
	} else {
		fmt.Println(".env not found, assuming CI environment")
	}

	// Start docker services
	fmt.Println("Starting docker compose services...")
	up := exec.Command("docker", "compose", "-f", "docker-compose.yaml", "up", "-d")
	up.Dir = "../../../"
	if output, err := up.CombinedOutput(); err != nil {
		fmt.Println("Docker startup failed:", string(output))
		os.Exit(1)
	}

	// Setup temp directory for JSON files
	tempDir, err := os.MkdirTemp("", "test_abilities")
	if err != nil {
		fmt.Println("Failed to create temp dir:", err)
		os.Exit(1)
	}
	os.Setenv("INITIAL_ABILITIES_FILE_PATH", filepath.Join(tempDir, "playersInitialAbilities.json"))
	os.Setenv("ABILITIES_FILE_PATH", filepath.Join(tempDir, "abilities.json"))
	fmt.Println("Setup temp JSON dir:", tempDir)

	// Run backend server in background
	go func() {
		_ = server.Start(server.ServerOptions{
			EnvFilePath: "../../../.env",
			Quiet:       true,
		})
	}()

	// Wait until the server is ready
	if err := waitForPort("localhost:3009", 10*time.Second); err != nil {
		fmt.Println("Server failed to start:", err)
		tearDownDocker(tempDir)
		os.Exit(1)
	}

	// Wait until the configurator server is ready
	if err := waitForPort("localhost:8080", 10*time.Second); err != nil {
		fmt.Println("Configurator server failed to start:", err)
		tearDownDocker(tempDir)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Always teardown
	tearDownDocker(tempDir)

	os.Exit(code)
}

func tearDownDocker(tempDir string) {
	fmt.Println("Tearing down docker compose services...")

	cmd := exec.Command("docker", "compose", "-f", "docker-compose.yaml", "down")
	cmd.Dir = "../../../"
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Docker teardown error:", err, string(output))
	}

	if tempDir != "" {
		_ = os.RemoveAll(tempDir)
	}
}

// ----------------------
// Test suite definition
// ----------------------

type ServerTestSuite struct {
	suite.Suite
	paladin   *TestClient
	barbarian *TestClient
	test      *TestClient
}

func (suite *ServerTestSuite) SetupSuite() {
	suite.paladin = connectClient("paladin")
	suite.barbarian = connectClient("barbarian")
}

func (suite *ServerTestSuite) TearDownSuite() {
	suite.paladin.Close()
	suite.barbarian.Close()
	// close test conn
}

// ----------------------
// Actual test
// ----------------------

func (suite *ServerTestSuite) Test_FullGameplayFlow() {
	suite.Run("DuplicateConnection", suite.testDuplicateCharacterConnection)
	suite.Run("Projectiles", suite.testProjectileDamagesTarget)
	suite.Run("PersistenceAfterReconnect", suite.testPersistenceAfterReconnect)
	suite.Run("ChangePlayersInitialAbilities", suite.testChangePlayersInitialAbilities)
	suite.Run("BuffDamageTarget", suite.testBuffDamageTarget)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) testDuplicateCharacterConnection() {
	url := "ws://localhost:3009/ws?character=barbarian"

	conn2, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(suite.T(), err, "second client should connect initially")
	defer conn2.Close()

	_, data, err := conn2.ReadMessage()
	assert.NoError(suite.T(), err, "should read close/error message")

	assert.Contains(suite.T(), string(data), "character already in use")
}

func (suite *ServerTestSuite) testProjectileDamagesTarget() {
	// Paladin casts projectile at barbarian
	err := suite.paladin.Send("ability_cast", map[string]interface{}{
		"id": "1", // hardcoded ability id, what to do here?
		"abilityParameters": map[string]interface{}{
			"TargetPosition": map[string]float64{
				"x": 5,
				"z": 0,
			},
		},
	})
	suite.Require().NoError(err)

	// Wait until the diff includes at least one projectile
	_, err = suite.paladin.WaitForGameStateDiff(3*time.Second, func(body map[string]interface{}) bool {
		projectiles, ok := body["projectiles"].([]interface{})
		return ok && len(projectiles) > 0
	})
	suite.Require().NoError(err, "projectile should appear")

	// Wait until barbarian HP drops
	_, err = suite.barbarian.WaitForGameStateDiff(5*time.Second, func(body map[string]interface{}) bool {
		players, ok := body["players"].([]interface{})
		if !ok {
			return false
		}
		for _, p := range players {
			pm := p.(map[string]interface{})
			if pm["id"] == "barbarian" {
				hp := pm["currentHealth"]
				if hpVal, ok := hp.(float64); ok && int(hpVal) < 100 {
					return true
				}
			}
		}
		return false
	})
	suite.Require().NoError(err, "barbarian HP should drop after hit")
}

func (suite *ServerTestSuite) testPersistenceAfterReconnect() {
	// Wait enough time so the state is saved to the db
	time.Sleep(5 * time.Second)

	suite.barbarian.Close()

	time.Sleep(500 * time.Millisecond)
	suite.barbarian = connectClient("barbarian")

	// Wait for initial diff with player state
	diff, err := suite.barbarian.WaitForGameStateDiff(3*time.Second, func(body map[string]interface{}) bool {
		players, ok := body["players"].([]interface{})
		if !ok {
			return false
		}
		for _, p := range players {
			pm := p.(map[string]interface{})
			if pm["id"] == "barbarian" {
				return true
			}
		}
		return false
	})
	suite.Require().NoError(err)

	players := diff["players"].([]interface{})
	for _, p := range players {
		pm := p.(map[string]interface{})
		if pm["id"] == "barbarian" {
			hp := pm["currentHealth"]
			if hpVal, ok := hp.(float64); ok {
				suite.True(hpVal < 100, "barbarian HP should remain reduced after reconnect")
			}
		}
	}
}

func (suite *ServerTestSuite) testChangePlayersInitialAbilities() {
	// avoid having to disconnect already connected characters because they interfere with the assert
	suite.barbarian.Close()
	suite.paladin.Close()

	playersInitialAbilitiesIds := [4]string{"1", "2", "3", "5"}
	jsonPayload, _ := json.Marshal(playersInitialAbilitiesIds)
	resp, err := http.Post("http://0.0.0.0:8080/playersInitialAbilities", "application/json", bytes.NewBuffer(jsonPayload))
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	suite.test = connectClient("test")

	_, err = suite.test.WaitForGameStateDiff(3*time.Second, func(body map[string]interface{}) bool {
		players, ok := body["players"].([]interface{})
		if !ok {
			return false
		}

		// Look for a character that has ability "5" (Buff Damage) which should be the new character
		return lo.ContainsBy(players, func(player interface{}) bool {
			return characterHasAbility(player, "5")
		})
	})
	suite.Require().NoError(err, "Should find a character with ability '5' (Buff Damage) after changing initial abilities")
}

func (suite *ServerTestSuite) testBuffDamageTarget() {
	// Wait for the test character to be fully connected and get initial stats
	initialStats, err := suite.test.WaitForGameStateDiff(3*time.Second, func(body map[string]interface{}) bool {
		players, ok := body["players"].([]interface{})
		if !ok {
			return false
		}
		return lo.ContainsBy(players, func(player interface{}) bool {
			pm := player.(map[string]interface{})
			if id, ok := pm["id"].(string); ok && id == suite.test.CharacterID {
				// Check if stats are present (means character is fully loaded)
				_, hasStats := pm["stats"]
				return hasStats
			}
			return false
		})
	})
	suite.Require().NoError(err, "Should receive initial character stats")

	// Extract initial damage stat
	var initialDamage int64
	players := initialStats["players"].([]interface{})
	for _, p := range players {
		pm := p.(map[string]interface{})
		if id, ok := pm["id"].(string); ok && id == suite.test.CharacterID {
			if stats, ok := pm["stats"].(map[string]interface{}); ok {
				if damage, ok := stats["damage"].(float64); ok {
					initialDamage = int64(damage)
					break
				}
			}
		}
	}
	suite.Require().True(initialDamage > 0, "Should have initial damage stat")

	// Test character casts buff on itself
	err = suite.test.Send("ability_cast", map[string]interface{}{
		"id": "5",
		"abilityParameters": map[string]interface{}{
			"TargetId": suite.test.CharacterID,
		},
	})
	suite.Require().NoError(err)

	// Wait for buff to be applied and verify damage stat increased
	buffedStats, err := suite.test.WaitForGameStateDiff(3*time.Second, func(body map[string]interface{}) bool {
		players, ok := body["players"].([]interface{})
		if !ok {
			return false
		}
		return lo.ContainsBy(players, func(player interface{}) bool {
			pm := player.(map[string]interface{})
			if id, ok := pm["id"].(string); ok && id == suite.test.CharacterID {
				if stats, ok := pm["stats"].(map[string]interface{}); ok {
					_, hasDamageStat := stats["damage"]
					return hasDamageStat
				}
			}
			return false
		})
	})
	suite.Require().NoError(err, "Should receive buffed character stats")

	// Verify the damage stat is higher
	var buffedDamage int64
	players = buffedStats["players"].([]interface{})
	for _, p := range players {
		pm := p.(map[string]interface{})
		if id, ok := pm["id"].(string); ok && id == suite.test.CharacterID {
			if stats, ok := pm["stats"].(map[string]interface{}); ok {
				if damage, ok := stats["damage"].(float64); ok {
					buffedDamage = int64(damage)
					break
				}
			}
		}
	}
	suite.True(buffedDamage > initialDamage, fmt.Sprintf("Damage should be higher after buff, original value: %d, new value: %d", initialDamage, buffedDamage))

	// Wait for buff to expire (buff duration is 5000ms = 5 seconds)
	time.Sleep(6 * time.Second)

	// Wait for buff to expire and verify damage stat returned to original
	_, err = suite.test.WaitForGameStateDiff(3*time.Second, func(body map[string]interface{}) bool {
		players, ok := body["players"].([]interface{})
		if !ok {
			return false
		}
		return lo.ContainsBy(players, func(player interface{}) bool {
			pm := player.(map[string]interface{})
			if id, ok := pm["id"].(string); ok && id == suite.test.CharacterID {
				if stats, ok := pm["stats"].(map[string]interface{}); ok {
					if damage, ok := stats["damage"].(float64); ok {
						// Check if damage is back to initial value
						return int64(damage) == initialDamage
					}
				}
			}
			return false
		})
	})
	suite.Require().NoError(err, "Damage should return to original value after buff expires")
}

// ----------------------
// Helper functions
// ----------------------

func characterHasAbility(c interface{}, abilityID string) bool {
	pm := c.(map[string]interface{})
	abilities, ok := pm["abilities"].([]interface{})
	if !ok {
		return false
	}
	return lo.ContainsBy(abilities, func(ability interface{}) bool {
		abilityMap, ok := ability.(map[string]interface{})
		if !ok {
			return false
		}
		id, ok := abilityMap["id"].(string)
		return ok && id == abilityID
	})
}

// ----------------------
// Helper: wait for port
// ----------------------

func waitForPort(address string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %s", address)
}

type WSMessage struct {
	ActionType string                 `json:"actionType"`
	Body       map[string]interface{} `json:"body"`
}

type TestClient struct {
	Conn        *websocket.Conn
	MsgChan     chan WSMessage
	CharacterID string
}

func connectClient(characterID string) *TestClient {
	url := fmt.Sprintf("ws://localhost:3009/ws?character=%s", characterID)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		panic(fmt.Sprintf("failed to connect %s: %v", characterID, err))
	}

	client := &TestClient{
		Conn:        conn,
		MsgChan:     make(chan WSMessage, 10),
		CharacterID: "", // Will be set from first Player message
	}

	// Reader goroutine
	go func() {
		defer close(client.MsgChan)
		for {
			_, buf, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var msg WSMessage
			if err := json.Unmarshal(buf, &msg); err == nil {
				// Capture character ID from first Player message
				if msg.ActionType == "Player" && client.CharacterID == "" {
					if id, ok := msg.Body["id"].(string); ok {
						client.CharacterID = id
					}
				}
				client.MsgChan <- msg
			}
		}
	}()

	return client
}

func (c *TestClient) Send(actionType string, body any) error {
	msg := map[string]interface{}{
		"actionType": actionType,
		"body":       body,
	}
	data, _ := json.Marshal(msg)
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}

func (c *TestClient) Close() { c.Conn.Close() }

func (c *TestClient) WaitForGameStateDiff(
	timeout time.Duration,
	match func(map[string]interface{}) bool,
) (map[string]interface{}, error) {

	deadline := time.After(timeout)
	for {
		select {
		case <-deadline:
			return nil, fmt.Errorf("timeout waiting for matching GameState diff")
		case msg, ok := <-c.MsgChan:
			if !ok {
				return nil, fmt.Errorf("connection closed before match")
			}
			if strings.ToLower(msg.ActionType) != "gamestate" {
				continue
			}
			if match(msg.Body) {
				return msg.Body, nil
			}
		}
	}
}
