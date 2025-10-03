package integration_test

import (
	"backend/api/dtos"
	"backend/connection"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func TestBasicWebSocketConnection(t *testing.T) {
	url := "ws://localhost:3009/ws?character=barbarian"

	conn, err := websocket.Dial(url, "", "http://localhost/")
	assert.NoError(t, err, "should connect to server")
	defer conn.Close()

	buf := make([]byte, 1024)
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	assert.NoError(t, err, "should receive message from server")

	var message connection.WebSocketMessage
	var character dtos.CharacterDTO
	err = json.Unmarshal(buf[:n], &message)
	assert.NoError(t, err, "should parse JSON correctly")

	assert.IsType(t, character, message.Body)
	character = message.Body.(dtos.CharacterDTO)

	var damageStat int64 = 20
	assert.Equal(t, "barbarian", character.Id)
	assert.Equal(t, (*character.Stats)["damage"], damageStat)
}
