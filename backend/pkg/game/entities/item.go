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

func NewItem(id, name string, mechanics ...Mechanic) *Item {
	item := &Item{
		id:   id,
		name: name,
	}

	if len(mechanics) > 0 {
		item.isConsumable = true
		item.mechanics = mechanics
	} else {
		item.isConsumable = false
		item.mechanics = nil
	}

	return item
}

func (item *Item) ExecuteMechanics(caster *Character, targetId string, gs *GameState) {
	for _, mechanic := range item.mechanics {
		if handler, exists := mechanicHandlers[mechanic.MechanicType]; exists {
			if err := handler(caster, targetId, gs, mechanic.Params); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("no handler found for effect type: %s/n", mechanic.MechanicType)
		}
	}
}
