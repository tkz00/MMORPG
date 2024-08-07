package game

type PlayerStats struct {
	maxHealth     int
	currentHealth int
}

func (stats PlayerStats) GetMaxHealth() int {
	return stats.maxHealth
}

func (stats PlayerStats) GetCurrentHealth() int {
	return stats.currentHealth
}
