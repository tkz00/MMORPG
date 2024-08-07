package dtos

type AbilityDTO struct {
	Id 					string	 	`json:"id"`
	Name 				string	 	`json:"name"`
	Range				float64		`json:"range"`
	Cooldown			int64		`json:"cooldown"`
	RemainingCooldown	float64		`json:"remainingCooldown"`
}
