package connection

import (
	"backend/api/dtos"
	"backend/pkg/utils"
	"encoding/json"
	"fmt"
)

type WebSocketMessage struct {
	Body       dtos.DTO `json:"body"`
	ActionType string   `json:"actionType"`
}

func CreateWebSocketResponse(body dtos.DTO) WebSocketMessage {
	return WebSocketMessage{
		Body:       body,
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
		Body       json.RawMessage `json:"body"`
		ActionType string          `json:"actionType"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	wr.ActionType = tmp.ActionType

	switch tmp.ActionType {
	case "position":
		var pos dtos.PositionDTO
		if err := json.Unmarshal(tmp.Body, &pos); err != nil {
			return err
		}
		wr.Body = pos
	case "ability_cast":
		var ability dtos.AbilityCastDTO
		if err := json.Unmarshal(tmp.Body, &ability); err != nil {
			return err
		}
		wr.Body = ability
	case "respawn":
	case "use_item":
		var useItem dtos.UseItemDTO
		if err := json.Unmarshal(tmp.Body, &useItem); err != nil {
			return err
		}
		wr.Body = useItem
	case "equip_item":
		var equipItem dtos.EquipItemDTO
		if err := json.Unmarshal(tmp.Body, &equipItem); err != nil {
			return err
		}
		wr.Body = equipItem
	default:
		return fmt.Errorf("unknown message type: %s", tmp.ActionType)
	}

	return nil
}
