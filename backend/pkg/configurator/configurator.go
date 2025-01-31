package configurator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"tkz00/backend/pkg/game/entities"
	"tkz00/backend/pkg/game/repository"

	"github.com/gin-gonic/gin"
)

const ABILITIES_FILE_NAME = "abilities.json"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/abilities", GetAbilities)
	r.GET("/abilities/:id", GetAbility)
	r.POST("/ability", CreateAbility)
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

func CreateAbility(c *gin.Context) {
	abilities, err := loadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	newAbility := ConfiguratorAbility{}
	if err := c.BindJSON(&newAbility); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if newAbility.Name == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("`name` required"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "`name` required",
		})
		return
	}

	if newAbility.Cooldown == nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("`cooldown` required"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "`cooldown` required",
		})
		return
	}

	if newAbility.RangeValue == nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("`range` required"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "`range` required",
		})
		return
	}

	if newAbility.Targeting == nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("`targeting` required"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "`targeting` required",
		})
		return
	}

	if newAbility.CharacterState == nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("`character_state` required"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "`character_state` required",
		})
		return
	}

	lastId := -1
	for idStr := range abilities {
		id, _ := strconv.Atoi(idStr)
		if id > lastId {
			lastId = id
		}
	}

	lastIdStr := strconv.Itoa(lastId + 1)
	newAbility.ID = lastIdStr
	abilities[lastIdStr] = newAbility
	saveAbilitiesToFile(abilities, ABILITIES_FILE_NAME)
	c.JSON(http.StatusAccepted, newAbility)
}

func UpdateAbility(c *gin.Context) {
	id := c.Param("id")
	abilities, err := loadAbilitiesFromFile(ABILITIES_FILE_NAME)
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	if ability, ok := abilities[id]; ok {
		updatedAbility := ConfiguratorAbility{}
		if err := c.BindJSON(&updatedAbility); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		fmt.Println(updatedAbility)

		if updatedAbility.Name != "" {
			ability.Name = updatedAbility.Name
		}

		if updatedAbility.Cooldown != nil {
			ability.Cooldown = updatedAbility.Cooldown
		}

		if updatedAbility.RangeValue != nil {
			ability.RangeValue = updatedAbility.RangeValue
		}

		if updatedAbility.Targeting != nil {
			ability.Targeting = updatedAbility.Targeting
		}

		if updatedAbility.CharacterState != nil {
			ability.CharacterState = updatedAbility.CharacterState
		}

		if len(updatedAbility.Mechanics) > 0 {
			ability.Mechanics = updatedAbility.Mechanics
		}

		abilities[id] = ability
		saveAbilitiesToFile(abilities, ABILITIES_FILE_NAME)
		c.JSON(http.StatusAccepted, ability)
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
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	RangeValue     *float64            `json:"range"`
	Cooldown       *int64              `json:"cooldown"`
	Targeting      *entities.Targeting `json:"targeting"`
	CharacterState *entities.Action    `json:"character_state"`
	Mechanics      []entities.Mechanic `json:"mechanics"`
}

// Convert function
func ConvertToConfiguratorAbility(ability entities.Ability) ConfiguratorAbility {
	targeting := entities.Targeting(
		ability.Targeting(),
	) // I don't know why I have to do this, but it doesn't work otherwise
	rangeValue := ability.Range()
	cooldown := ability.Cooldown()
	characterState := ability.CharacterState()
	return ConfiguratorAbility{
		ID:             ability.Id(),
		Name:           ability.Name(),
		RangeValue:     &rangeValue,
		Cooldown:       &cooldown,
		Targeting:      &targeting,
		CharacterState: &characterState,
		Mechanics:      ability.Mechanics(),
	}
}
