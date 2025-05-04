package configurator

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/abilities", GetAbilities)
	r.GET("/abilities/:id", GetAbility)
	r.POST("/ability", CreateAbility)
	r.PATCH("/ability/:id", UpdateAbility)
	r.DELETE("/ability/:id", DeleteAbility)
	return r
}

func RunSeeds() {
	abilities := GetSeedsAbilities()
	configuratorAbilities := make(map[string]ConfiguratorAbility, len(abilities))
	for i, ability := range abilities {
		configuratorAbilities[i] = ConvertToConfiguratorAbility(*ability)
	}

	SaveAbilitiesToFile(configuratorAbilities)
	playerAbilitiesIds := lo.Keys(configuratorAbilities)
	sort.Strings(playerAbilitiesIds)
	SavePlayerInitialAbilities(playerAbilitiesIds[1:5])
}

func Run() {
	r := SetupRouter()
	fmt.Println("Starting server on port 8080...")
	r.Run()
}

func GetAbilities(c *gin.Context) {
	abilities, err := LoadAbilitiesFromFile()
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
	abilities, err := LoadAbilitiesFromFile()
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
	var newAbility ConfiguratorAbility
	if err := c.ShouldBindJSON(&newAbility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateAbility(newAbility); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	abilities, err := LoadAbilitiesFromFile()
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	newID := getNextAbilityID(abilities)
	newAbility.ID = newID
	abilities[newID] = newAbility
	SaveAbilitiesToFile(abilities)
	c.JSON(http.StatusCreated, newAbility)
}

func UpdateAbility(c *gin.Context) {
	id := c.Param("id")
	abilities, err := LoadAbilitiesFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load abilities"})
		return
	}
	ability, exists := abilities[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
		return
	}

	updatedAbility := ConfiguratorAbility{}
	if err := c.BindJSON(&updatedAbility); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

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
	SaveAbilitiesToFile(abilities)
	c.JSON(http.StatusAccepted, ability)
}

func DeleteAbility(c *gin.Context) {
	id := c.Param("id")
	abilities, err := LoadAbilitiesFromFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load abilities"})
		return
	}

	if _, exists := abilities[id]; exists {
		delete(abilities, id)
		SaveAbilitiesToFile(abilities)
		c.JSON(http.StatusOK, gin.H{"message": "ability deleted"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
}

func validateAbility(a ConfiguratorAbility) error {
	if a.Name == "" {
		return errors.New("`name` required")
	}
	if a.Cooldown == nil {
		return errors.New("`cooldown` required")
	}
	if a.RangeValue == nil {
		return errors.New("`range` required")
	}
	if a.Targeting == nil {
		return errors.New("`targeting` required")
	}
	if a.CharacterState == nil {
		return errors.New("`character_state` required")
	}
	return nil
}

func getNextAbilityID(abilities map[string]ConfiguratorAbility) string {
	maxID := 0
	for id := range abilities {
		if num, err := strconv.Atoi(id); err == nil && num > maxID {
			maxID = num
		}
	}
	return strconv.Itoa(maxID + 1)
}
