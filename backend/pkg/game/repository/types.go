package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StatsMap map[string]int64

// Value implements driver.Valuer
func (s StatsMap) Value() (driver.Value, error) {
	return json.Marshal(s) // store as JSONB
}

// Scan implements sql.Scanner
func (s *StatsMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StatsMap: %v", value)
	}
	return json.Unmarshal(b, s) // load into map[string]int64
}

type ItemsMap map[string]int64

// Value implements driver.Valuer
func (s ItemsMap) Value() (driver.Value, error) {
	return json.Marshal(s) // store as JSONB
}

// Scan implements sql.Scanner
func (s *ItemsMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ItemsMap: %v", value)
	}
	return json.Unmarshal(b, s) // load into map[string]int64
}

type EquipmentMap map[string]string

// Value implements driver.Valuer
func (s EquipmentMap) Value() (driver.Value, error) {
	return json.Marshal(s) // store as JSONB
}

// Scan implements sql.Scanner
func (s *EquipmentMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan EquipmentMap: %v", value)
	}
	return json.Unmarshal(b, s) // load into map[string]int64
}
