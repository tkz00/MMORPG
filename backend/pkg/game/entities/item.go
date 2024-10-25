package entities

import "fmt"

type Item struct {
	id           string
	name         string
	isConsumable bool
	mechanics    []Mechanic
}

func (itemTemplate Item) Id() string {
	return itemTemplate.id
}

func (itemTemplate Item) Name() string {
	return itemTemplate.name
}

type Mechanic struct {
	MechanicType string                 // Type of effect (e.g., "heal", "unlock")
	Params       map[string]interface{} // Dynamic parameters for the effect (e.g., healing amount, target)
}

type MechanicHandler func(c *Character, params map[string]interface{}) error

var mechanicHandlers = map[string]MechanicHandler{}

// I don't know where this should be
func RegisterMechanicHandler(mechanicType string, handler MechanicHandler) {
	mechanicHandlers[mechanicType] = handler
}

func HealMechanic(c *Character, params map[string]interface{}) error {
	if amount, ok := params["amount"].(int); ok {
		c.HealthVariation(amount)
		return nil
	}
	return fmt.Errorf("missing or invalid 'amount' parameter")
}

func DamageMechanic(c *Character, params map[string]interface{}) error {
	if amount, ok := params["amount"].(int); ok {
		c.HealthVariation(-amount)
		return nil
	}
	return fmt.Errorf("missing or invalid 'amount' parameter")
}
