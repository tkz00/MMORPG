package connection

import (
	"encoding/json"
	"fmt"
	"unnamed-mmo/backend/types"
	"unnamed-mmo/backend/utils"
)

type WebSocketMessage struct {
	Body types.DTO            `json:"body"`
	ActionType string   `json:"actionType"`
}

func CreateWebSocketResponse(body types.DTO) WebSocketMessage {
	return WebSocketMessage{
		Body: body,
		ActionType: body.GetType(),
	}
}

func (wsr WebSocketMessage) Serialize() []byte {
	data, err := utils.ToJSON(wsr)

	if err != nil {
		fmt.Println("Error serializing "+wsr.ActionType, err)
		return nil
	}

	return data
}

func (wr *WebSocketMessage) UnmarshalJSON(data []byte) error {
    var tmp struct {
        Body json.RawMessage    `json:"body"`
        ActionType string       `json:"actionType"`
    }

    if err := json.Unmarshal(data, &tmp); err != nil {
        return err
    }

    wr.ActionType = tmp.ActionType

    switch tmp.ActionType {
    case "Position":
        var pos types.PositionDTO
        if err := json.Unmarshal(tmp.Body, &pos); err != nil {
            return err
        }
        wr.Body = pos
    case "AbilityCast":
        var ability types.AbilityCastDTO
        if err := json.Unmarshal(tmp.Body, &ability); err != nil {
            return err
        }
        wr.Body = ability
    default:
        return fmt.Errorf("unknown message type: %s", tmp.ActionType)
    }

    return nil
}
