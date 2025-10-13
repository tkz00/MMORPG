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

type ServerTestSuite struct {
	suite.Suite
}

func (suite *ServerTestSuite) SetupSuite() {
	// Setup test env var
	os.Setenv("GO_ENV", "test")

	// Load env file by absolute path dynamically
	if os.Getenv("GO_ENV") == "test" {
		cwd, _ := os.Getwd()
		root := filepath.Join(cwd, "..", "..", "..")
		if err := godotenv.Load(filepath.Join(root, ".env")); err != nil {
			fmt.Println("Warning: failed to load .env:", err)
		}
	}

	// Run postgresql docker
	cmd := exec.Command("docker", "compose", "-f", "docker-compose.yaml", "up", "-d")
	cmd.Dir = "../../../"
	suite.T().Log("Starting docker compose services...")
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err, string(output))

	go setupServer()
	err = waitForPort("localhost:3009", 10*time.Second)
	suite.Require().NoError(err, "backend server didn't start in time")
}

func (suite *ServerTestSuite) TearDownSuite() {
	cmd := exec.Command("docker", "compose", "-f", "../../../docker-compose.yaml", "down")
	suite.T().Log("Tearing down docker compose services...")
	output, err := cmd.CombinedOutput()
	suite.Require().NoError(err, string(output))
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

	const PORT string = "3009"
	go connection.NewServer().StartConnection(PORT)

	select {}
}

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
