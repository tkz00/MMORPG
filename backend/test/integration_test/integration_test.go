package integration_test

import (
	"backend/api/dtos"
	"backend/connection"
	"backend/pkg/configurator"
	"backend/pkg/game/repository"
	"backend/pkg/handlers"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	go setupServer()

	// Wait until the server is ready
	if err := waitForPort("localhost:3009", 10*time.Second); err != nil {
		fmt.Println("Server failed to start:", err)
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
}

func setupServer() {
	if os.Getenv("GO_ENV") == "test" {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}

	configurator.RunSeeds()

	go configurator.Run()

	if err := repository.ConnectPostgres(); err != nil {
		panic(err)
	}

	repository.RunSeeds()
	handlers.RegisterRoutes()

	const PORT = "3009"
	go connection.NewServer().StartConnection(PORT)

	select {} // keep goroutine alive
}

// ----------------------
// Actual test
// ----------------------

func (suite *ServerTestSuite) TestBasicWebSocketConnection() {
	url := "ws://localhost:3009/ws?character=barbarian"

	conn, err := websocket.Dial(url, "", "http://localhost/")
	assert.NoError(suite.T(), err, "should connect to server")
	defer conn.Close()

	buf := make([]byte, 1024)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	assert.NoError(suite.T(), err, "should receive message from server")

	var message connection.WebSocketMessage
	var character dtos.CharacterDTO
	err = json.Unmarshal(buf[:n], &message)
	assert.NoError(suite.T(), err, "should parse JSON correctly")

	assert.IsType(suite.T(), character, message.Body)
	character = message.Body.(dtos.CharacterDTO)

	var damageStat int64 = 20
	assert.Equal(suite.T(), "barbarian", character.Id)
	assert.Equal(suite.T(), (*character.Stats)["damage"], damageStat)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
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
