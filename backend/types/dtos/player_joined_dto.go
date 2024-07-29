package dtos

import "unnamed-mmo/backend/types"

type PlayerJoinedDTO struct {
	Id			string			`json:"id"`
	Abilities	[]AbilityDTO	`json:"abilities"`
}

func (m Mapper) PlayerToJoinedDTO(player types.Player) *PlayerJoinedDTO {
	playerAbilities := make([]AbilityDTO, 0)
	for _, ability := range player.GetAbilities() {
		playerAbilities = append(playerAbilities, AbilityToDTO(ability))
	}

	return &PlayerJoinedDTO {
		Id: player.GetId(),
		Abilities: playerAbilities,
	}
}

func (p PlayerJoinedDTO) GetType() string {
	return "PlayerJoinedDTO"
}
