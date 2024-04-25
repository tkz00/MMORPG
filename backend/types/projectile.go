package types

type Projectile struct {
	direction Position
	stats     PlayerStats
}

func CreateProjectile(direction Position) *Projectile {
	return &Projectile{
		direction: direction,
		stats:     PlayerStats{},
	}
}