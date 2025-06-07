// I don't know if this should be in entities, prolly doesn't matter, same for item
package entities

import (
	"fmt"
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

func (inv *Inventory) EquipItem(itemId string) error {
	item := existingItems[itemId]
	if equip, ok := item.(*Equipment); ok {
		inv.equipped[equip.equipmentType] = equip
		return nil
	}
	return fmt.Errorf("item %s is not equipment", itemId)
}

func (inv *Inventory) UnequipItem(equipmentType EquipmentType) {
	delete(inv.equipped, equipmentType)
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
}
