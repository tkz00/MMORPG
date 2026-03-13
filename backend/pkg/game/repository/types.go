package repository

import (
	"backend/pkg/game/entities"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StatsMap map[string]int64

func (s StatsMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StatsMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StatsMap: %v", value)
	}
	return json.Unmarshal(b, s)
}

type ItemsMap map[string]int64

func (s ItemsMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *ItemsMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ItemsMap: %v", value)
	}
	return json.Unmarshal(b, s)
}

type EquipmentMap map[string]string

func (s EquipmentMap) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *EquipmentMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan EquipmentMap: %v", value)
	}
	return json.Unmarshal(b, s)
}

type MechanicsList []entities.Mechanic

func (m MechanicsList) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MechanicsList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan MechanicsList: %v", value)
	}
	return json.Unmarshal(b, m)
}

type PerkStateEntry struct {
	StatName     string  `json:"stat_name"`
	FlatValue    int64   `json:"flat_value"`
	PercentValue float64 `json:"percent_value"`
}

type PerkStateList []PerkStateEntry

func (p PerkStateList) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PerkStateList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan PerkStateList: %v", value)
	}
	return json.Unmarshal(b, p)
}
