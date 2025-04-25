// I don't know if this should be in entities, prolly doesn't matter, same for item
package entities

import "fmt"

type ItemChange struct {
	Id       string `json:"id"`       // Template ID of the item
	Quantity int    `json:"quantity"` // New quantity of the item
}

type Inventory struct {
	items     map[*Item]int
	changeLog []ItemChange
}

func NewInventory() *Inventory {
	return &Inventory{
		items: make(map[*Item]int),
	}
}

func (looter *Inventory) Loot(loot *Inventory) {
	for item := range loot.items {
		looter.AddItem(item, loot.items[item])
	}
}

func (looter *Inventory) AddItem(item *Item, quantity int) {
	for itemInInventory := range looter.items {
		if itemInInventory.id == item.id {
			looter.items[item] += quantity
			if looter.items[item] == 0 {
				delete(looter.items, item)
			}
			looter.changeLog = append(looter.changeLog, ItemChange{
				Id:       item.id,
				Quantity: quantity,
			})
			return
		}
	}

	looter.items[item] = quantity
	looter.changeLog = append(looter.changeLog, ItemChange{
		Id:       item.id,
		Quantity: quantity,
	})
}

func (inv Inventory) CanConsume(itemId string) bool {
	for item := range inv.items {
		if item.id == itemId {
			return item.isConsumable
		}
	}
	return false
}

func (inv Inventory) GetItem(itemId string) *Item {
	for item := range inv.items {
		if item.id == itemId {
			return item
		}
	}
	return nil
}

func (inventory *Inventory) ChangeLogs() []ItemChange {
	changes := inventory.changeLog
	inventory.changeLog = []ItemChange{}
	return changes
}

func (inventory *Inventory) GetInventory() []ItemChange {
	invSlice := []ItemChange{}
	for item, qty := range inventory.items {
		invSlice = append(invSlice, ItemChange{Id: item.id, Quantity: qty})
	}
	return invSlice
}

func (inv *Inventory) PrintInventory() {
	if len(inv.items) > 0 {
		for item, quantity := range inv.items {
			fmt.Printf("%s, %d - ", item.name, quantity)
		}
		fmt.Print("\n")
	}
}
