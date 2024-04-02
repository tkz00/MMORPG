package types

type PlayerDTO struct {
	Id            string      `json:"id"`
	MaxHealth     int         `json:"maxHealth"`
	CurrentHealth int         `json:"currentHealth"`
	Position      PositionDTO `json:"position"`
}
