package types

import (
	"fmt"
	"unnamed-mmo/backend/utils"
)

type WebSocketResponse struct {
	Body DTO    `json:"body"`
	Type string `json:"type"`
	Code string `json:"code"`
}

func CreateWebSocketResponse(body DTO, code string) WebSocketResponse {
	return WebSocketResponse{
		Body: body,
		Type: body.GetType(),
		Code: code,
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