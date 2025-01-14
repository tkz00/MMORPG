package configurator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/gin-gonic/gin"
)

const ABILITIES_FILE_NAME = "abilities.json"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/abilities", GetAbilities)
	r.GET("/abilities/:id", GetAbility)
	r.PATCH("/ability/:id", UpdateAbility)
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

func GetAbilities(c *gin.Context) {
	abilities, err := loadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"abilities": abilities,
	})
}

func GetAbility(c *gin.Context) {
	id := c.Param("id")
	abilities, err := loadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	if ability, ok := abilities[id]; ok {
		c.JSON(http.StatusOK, ability)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "ability not found"})
}

func UpdateAbility(c *gin.Context) {
	id := c.Param("id")
	abilities, err := loadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	if ability, ok := abilities[id]; ok {
		fmt.Println(ability)
		updatedAbility := ConfiguratorAbility{}
		if err := c.BindJSON(&updatedAbility); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if updatedAbility.Name != "" {
			ability.Name = updatedAbility.Name
		}

		if updatedAbility.Cooldown != 0 {
			ability.Cooldown = updatedAbility.Cooldown
		}

		if updatedAbility.RangeValue != 0 {
			ability.RangeValue = updatedAbility.RangeValue
		}

		if updatedAbility.Targeting != nil {
			ability.Targeting = updatedAbility.Targeting
		}

		abilities[id] = ability
		saveAbilitiesToFile(abilities, ABILITIES_FILE_NAME)
		c.JSON(http.StatusAccepted, &ability)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "ability not found"})
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
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	RangeValue float64             `json:"rangeValue"`
	Cooldown   int64               `json:"cooldown"`
	Targeting  *entities.Targeting `json:"targeting"`
	// CharacterState entities.Action     `json:"characterState"`
	// Mechanics []entities.Mechanic `json:"mechanics"`
}

// Convert function
func ConvertToConfiguratorAbility(ability entities.Ability) ConfiguratorAbility {
	targeting := entities.Targeting(
		ability.Targeting(),
	) // I don't know why I have to do this, but it doesn't work otherwise
	return ConfiguratorAbility{
		ID:         ability.Id(),
		Name:       ability.Name(),
		RangeValue: ability.Range(),
		Cooldown:   ability.Cooldown(),
		Targeting:  &targeting,
		// CharacterState: ability.characterState,
		// Mechanics: ability.mechanics,
	}
}
