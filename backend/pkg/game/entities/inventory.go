// I don't know if this should be in entities, prolly doesn't matter, same for item
package entities

import (
	"fmt"
)

type Inventory struct {
	items map[string]int64
}

func NewInventory() *Inventory {
	return &Inventory{
		items: make(map[string]int64),
	}
}

func (looter *Inventory) Loot(loot *Inventory) {
	for item, qty := range loot.items {
		looter.AddItem(item, qty)
	}
}

func (looter *Inventory) AddItem(itemId string, quantity int64) {
	for itemInInventoryId := range looter.items {
		if itemInInventoryId == itemId {
			looter.items[itemId] += quantity
			if looter.items[itemId] == 0 {
				delete(looter.items, itemId)
			}
			return
		}
	}

	looter.items[itemId] = quantity
}

func (inv Inventory) CanConsume(itemId string) bool {
	for itemInInventoryId := range inv.items {
		if itemInInventoryId == itemId {
			return GetItem(itemId).isConsumable
		}
	}
	return false
}

func GetItem(itemId string) *Item {
	return existingItems[itemId]
}

func (inv Inventory) GetInventory() map[string]int64 {
	return inv.items
}

func (inv *Inventory) PrintInventory() {
	if len(inv.items) > 0 {
		for itemId, quantity := range inv.items {
			fmt.Printf("%s, %d - ", itemId, quantity)
		}
		fmt.Print("\n")
	}
}

var existingItems = map[string]*Item{
	"0": NewItem(
		"0",
		"small health potion",
		Mechanic{
			MechanicType: "heal",
			Params:       map[string]interface{}{"base_amount": 40.0},
		},
	),
	"1": NewItem("1", "leather"),
}
