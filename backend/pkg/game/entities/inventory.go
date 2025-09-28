// I don't know if this should be in entities, prolly doesn't matter, same for item
package entities

import (
	"fmt"

	"github.com/samber/lo"
)

type Inventory struct {
	items    map[string]int64
	equipped map[EquipmentType]*Equipment
}

func NewInventory() *Inventory {
	return &Inventory{
		items:    make(map[string]int64),
		equipped: make(map[EquipmentType]*Equipment),
	}
}

func (looter *Inventory) Loot(loot *Inventory) {
	for item, qty := range loot.items {
		looter.AddItem(item, qty)
	}
	loot.Clear()
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
	if item, ok := existingItems[itemId].(*Item); ok {
		return item
	}
	if equip, ok := existingItems[itemId].(*Equipment); ok {
		return equip.Item
	}
	return nil
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

func (i *Inventory) Clear() {
	for et := range i.equipped {
		delete(i.equipped, et)
	}
	for k := range i.items {
		delete(i.items, k)
	}
}

func (inv *Inventory) EquipItem(itemId string) error {
	if lo.Contains(lo.Keys(inv.items), itemId) {
		item := existingItems[itemId]
		if equip, ok := item.(*Equipment); ok {
			inv.equipped[equip.equipmentType] = equip
			return nil
		}
		return fmt.Errorf("item %s is not equipment", itemId)
	}
	return fmt.Errorf("player doesn't hold item %s", itemId)
}

func (inv *Inventory) UnequipItem(itemId string) {
	for equipType, equip := range inv.equipped {
		if equip.Item.Id() == itemId {
			delete(inv.equipped, equipType)
			return
		}
	}
}

func (inv *Inventory) GetEquipped() map[EquipmentType]*Equipment {
	return inv.equipped
}

var existingItems = map[string]interface{}{
	"0": NewItem(
		"0",
		"small health potion",
		Mechanic{
			MechanicType: "heal",
			Params:       map[string]interface{}{"base_amount": 40.0},
		},
	),
	"1": NewItem("1", "leather"),
	"helm_001": NewEquipment(
		"helm_001",
		"leather helmet",
		Helmet,
		map[string]int64{"defense": 10},
	),
	"helm_002": NewEquipment(
		"helm_002",
		"steel helmet",
		Helmet,
		map[string]int64{"defense": 20},
	),
	"chest_001": NewEquipment(
		"chest_001",
		"leather armour",
		Chest,
		map[string]int64{"defense": 12},
	),
}
