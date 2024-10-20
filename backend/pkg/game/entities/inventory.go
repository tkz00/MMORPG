// I don't know if this should be in entities, prolly doesn't matter, same for item
package entities

import "fmt"

type ItemChange struct {
	Id       string `json:"id"`       // Template ID of the item
	Quantity int    `json:"quantity"` // New quantity of the item
}

type Inventory struct {
	items     []*Item
	changeLog []ItemChange
}

func (looter *Inventory) Loot(loot *Inventory) {
	for _, item := range loot.items {
		looter.AddItem(item)
	}
}

func (looter *Inventory) AddItem(item *Item) {
	for _, itemInInventory := range looter.items {
		if itemInInventory.template.id == item.template.id {
			itemInInventory.quantity += item.quantity
			// Log an update change
			looter.changeLog = append(looter.changeLog, ItemChange{
				Id:       item.template.id,
				Quantity: item.quantity,
			})
			return
		}
	}

	// Add a new item
	looter.items = append(looter.items, item)
	// Log an add change
	looter.changeLog = append(looter.changeLog, ItemChange{
		Id:       item.template.id,
		Quantity: item.quantity,
	})
}

func (inventory *Inventory) ChangeLogs() []ItemChange {
	// Return the list of changes since last check
	changes := inventory.changeLog
	// Clear the change log after sending it
	inventory.changeLog = []ItemChange{}
	return changes
}

func (inv *Inventory) PrintInventory() {
	if len(inv.items) > 0 {
		for _, item := range inv.items {
			fmt.Printf("%s, %d - ", item.template.name, item.quantity)
		}
		fmt.Print("\n")
	}
}
