package configurator

import (
	"encoding/json"
	"fmt"
	"os"
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/gin-gonic/gin"
)

const ABILITIES_FILE_NAME = "abilities.json"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/abilities", func(c *gin.Context) {
		abilities, err := loadAbilitiesFromFile(ABILITIES_FILE_NAME)
		if err != nil {
			fmt.Printf("Error loading abilities: %v\n", err)
			return
		}

		c.JSON(200, gin.H{
			"abilities": abilities,
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
	abilities := repository.GetPlayerAbilities()
	configuratorAbilities := make(map[string]ConfiguratorAbility, len(abilities))
	for i, ability := range abilities {
		configuratorAbilities[i] = ConvertToConfiguratorAbility(*ability)
	}

	saveAbilitiesToFile(configuratorAbilities, ABILITIES_FILE_NAME)

	r := SetupRouter()
	fmt.Println("Starting server on port 8080...")
	r.Run()
}

func saveAbilitiesToFile(abilities map[string]ConfiguratorAbility, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(abilities)
	if err != nil {
		return fmt.Errorf("error encoding abilities to JSON: %v", err)
	}

	return nil
}

func loadAbilitiesFromFile(filename string) (map[string]ConfiguratorAbility, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	abilities := make(map[string]ConfiguratorAbility)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&abilities)
	if err != nil {
		return nil, fmt.Errorf("error decoding abilities from JSON: %v", err)
	}

	return abilities, nil
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
