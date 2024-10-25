package entities

type Item struct {
	id           string
	name         string
	isConsumable bool
}

func (itemTemplate Item) Id() string {
	return itemTemplate.id
}

func (itemTemplate Item) Name() string {
	return itemTemplate.name
}
