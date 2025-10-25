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
	"golang.org/x/net/websocket"
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
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) testDuplicateCharacterConnection() {
	url := "ws://localhost:3009/ws?character=barbarian"

	conn2, err := websocket.Dial(url, "", "http://localhost/")
	assert.NoError(suite.T(), err, "second client should connect initially")
	defer conn2.Close()

	buf := make([]byte, 512)
	n, err := conn2.Read(buf)
	assert.NoError(suite.T(), err, "should read close/error message")

	assert.Contains(suite.T(), string(buf[:n]), "character already in use")
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
	Conn    *websocket.Conn
	MsgChan chan WSMessage
}

func connectClient(characterID string) *TestClient {
	url := fmt.Sprintf("ws://localhost:3009/ws?character=%s", characterID)
	conn, err := websocket.Dial(url, "", "http://localhost/")
	if err != nil {
		panic(fmt.Sprintf("failed to connect %s: %v", characterID, err))
	}

	client := &TestClient{
		Conn:    conn,
		MsgChan: make(chan WSMessage, 10),
	}

	// Reader goroutine
	go func() {
		defer close(client.MsgChan)
		buf := make([]byte, 8192)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			var msg WSMessage
			if err := json.Unmarshal(buf[:n], &msg); err == nil {
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
	_, err := c.Conn.Write(data)
	return err
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
