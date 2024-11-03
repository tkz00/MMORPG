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

func (gs GameState) ExecuteMechanics(mechanics []Mechanic, caster Character) GameState {
	for _, mechanic := range mechanics {
		handler, exists := mechanicHandlers[mechanic.MechanicType]
		if !exists {
			fmt.Printf("no handler found for effect type: %s/n", mechanic.MechanicType)
		}

		var err error
		gs, err = handler(caster, gs, mechanic.Params)

		if err != nil {
			fmt.Println(err)
		}
	}

	return gs
}
