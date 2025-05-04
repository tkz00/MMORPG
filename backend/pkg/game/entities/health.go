package entities

type Health struct {
	maxHealth     int
	currentHealth int
}

func NewHealth(totalHealth int) Health {
	return Health{maxHealth: totalHealth, currentHealth: totalHealth}
}

func (health Health) GetMaxHealth() int {
	return health.maxHealth
}

func (health Health) GetCurrentHealth() int {
	return health.currentHealth
}

func (h *Health) HealthVariation(variation int) {
	newHealth := h.GetCurrentHealth() + variation
	if newHealth > 0 {
		if h.GetMaxHealth() > newHealth {
			h.currentHealth += variation
		} else {
			h.currentHealth = h.maxHealth
		}
	} else {
		h.currentHealth = 0
	}
}

func (h Health) IsAlive() bool {
	return h.currentHealth > 0
}
