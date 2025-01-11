package configurator

import (
	"fmt"
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// Register your routes here
	r.GET("/abilities", func(c *gin.Context) {
		abilities := repository.GetPlayerAbilities()
		configuratorAbilities := make(map[string]ConfiguratorAbility, len(abilities))
		for i, ability := range abilities {
			configuratorAbilities[i] = ConvertToConfiguratorAbility(*ability)
		}

		c.JSON(200, gin.H{
			"abilities": configuratorAbilities,
		})
	})

	r.PATCH("/ability/:id", func(c *gin.Context) {
		id := c.Param("id")
		fmt.Println(id)
		c.JSON(200, id)
	})
	return r
}

func Run() {
	r := SetupRouter()
	fmt.Println("Starting server on port 8080...")
	r.Run() // Default port
}

// ConfiguratorAbility struct for API response
type ConfiguratorAbility struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	RangeValue float64            `json:"rangeValue"`
	Cooldown   int64              `json:"cooldown"`
	Targeting  entities.Targeting `json:"targeting"`
	// CharacterState entities.Action     `json:"characterState"`
	// Mechanics []entities.Mechanic `json:"mechanics"`
}

// Convert function
func ConvertToConfiguratorAbility(ability entities.Ability) ConfiguratorAbility {
	return ConfiguratorAbility{
		ID:         ability.Id(),
		Name:       ability.Name(),
		RangeValue: ability.Range(),
		Cooldown:   ability.Cooldown(),
		Targeting:  ability.Targeting(),
		// CharacterState: ability.characterState,
		// Mechanics: ability.mechanics,
	}
}
