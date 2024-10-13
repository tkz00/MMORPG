package entities

type Item struct {
	template ItemTemplate
	quantity int
}

func (item Item) Template() ItemTemplate {
	return item.template
}

func (item Item) Quantity() int {
	return item.quantity
}
