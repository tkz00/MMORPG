package entities

type ItemTemplate struct {
	id   string
	name string
}

func (itemTemplate ItemTemplate) Id() string {
	return itemTemplate.id
}

func (itemTemplate ItemTemplate) Name() string {
	return itemTemplate.name
}
