package repository

import (
	"backend/pkg/game/entities"

	"gorm.io/gorm"
)

type EffectDB struct {
	Id         string `gorm:"primaryKey"`
	EffectName string
	Trigger    string
	Mechanics  MechanicsList `gorm:"type:jsonb"`

	gorm.Model
}

func (EffectDB) TableName() string {
	return "effects"
}

func GetEffectsByIds(effectIDs []string) []entities.Effect {
	if len(effectIDs) == 0 {
		return []entities.Effect{}
	}

	var dbEffects []EffectDB
	err := DB.Where("id IN ?", effectIDs).Find(&dbEffects).Error
	if err != nil {
		// depending on your style:
		// return empty, or panic, or log
		return []entities.Effect{}
	}

	// Convert DB models → runtime effect objects
	effects := make([]entities.Effect, 0, len(dbEffects))
	for _, dbE := range dbEffects {
		effects = append(effects, convertEffectDBToEntity(dbE))
	}

	return effects
}

func convertEffectDBToEntity(e EffectDB) entities.Effect {
	return *entities.CreateEffect(e.Id,
		e.EffectName,
		entities.EffectTrigger(e.Trigger),
		convertMechanics(e.Mechanics),
	)
}

func convertMechanics(list MechanicsList) []entities.Mechanic {
	result := make([]entities.Mechanic, 0, len(list))

	for _, m := range list {
		result = append(result, entities.Mechanic{
			MechanicType: m.MechanicType,
			Params:       m.Params, // already map[string]interface{}
		})
	}

	return result
}

func GetEffectById(effectId string) (*entities.Effect, error) {
	var repoEffect EffectDB
	if err := DB.Where("id = ?", effectId).First(&repoEffect).Error; err != nil {
		return nil, err
	}

	effect := convertEffectDBToEntity(repoEffect)

	return &effect, nil
}

func SaveEffect(e *entities.Effect) error {
	dbEffect, err := EffectFromEntity(e)
	if err != nil {
		return err
	}

	return DB.Save(dbEffect).Error
}

func EffectFromEntity(e *entities.Effect) (*EffectDB, error) {
	return &EffectDB{
		Id:         e.Id(),
		EffectName: e.Name(),
		Trigger:    string(e.Trigger()),
		Mechanics:  MechanicsList(e.Mechanics()),
	}, nil
}
