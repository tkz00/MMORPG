package entities

type EquipmentType string

const (
	Helmet EquipmentType = "helmet"
	Chest  EquipmentType = "chest"
	Boots  EquipmentType = "boots"
)

type Equipment struct {
	*Item
	equipmentType EquipmentType
	stats         map[string]int64
}

func NewEquipment(id, name string, equipmentType EquipmentType, stats map[string]int64) *Equipment {
	return &Equipment{
		Item:          NewItem(id, name),
		equipmentType: equipmentType,
		stats:         stats,
	}
}

func (e *Equipment) GetEquipmentType() EquipmentType {
	return e.equipmentType
}

func (e *Equipment) GetStats() map[string]int64 {
	return e.stats
}

func (e *Equipment) SetStats(stats map[string]int64) {
	e.stats = stats
}
