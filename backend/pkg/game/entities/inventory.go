// I don't know if this should be in entities, prolly doesn't matter, same for item
package entities

import "fmt"

type Inventory struct {
	items []*Item
}

func (looter *Inventory) Loot(loot Inventory) {
	for _, item := range loot.items {
		looter.AddItem(item)
	}
}

func (looter *Inventory) AddItem(item *Item) {
	for _, itemInInventory := range looter.items {
		if itemInInventory.template.id == item.template.id {
			itemInInventory.quantity += item.quantity
			return
		}
	}

	looter.items = append(looter.items, item)
}

func (inv *Inventory) PrintInventory() {
	if len(inv.items) > 0 {
		for _, item := range inv.items {
			fmt.Printf("%s, %d - ", item.template.name, item.quantity)
		}
		fmt.Print("\n")
	}
}
