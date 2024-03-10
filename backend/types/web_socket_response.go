package types

import (
	"fmt"
	"unnamed-mmo/backend/utils"
)

type WebSocketResponse struct {
	Body DTO    `json:"body"`
	Type string `json:"type"`
}

func CreateWebSocketResponse(body DTO) WebSocketResponse {
	return WebSocketResponse{
		Body: body,
		Type: body.GetType(),
	}
}

func (wsr WebSocketResponse) Serialize() []byte {
	data, err := utils.ToJSON(wsr)

	if err != nil {
		fmt.Println("Error serializing "+wsr.Type, err)
		return nil
	}

	return data
}