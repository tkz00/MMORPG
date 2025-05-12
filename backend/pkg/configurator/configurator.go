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
	r.GET("/playersInitialAbilities", GetPlayersInitialAbilities)
	r.POST("/playersInitialAbilities", SetPlayersInitialAbilities)
	return r
}

func RunSeeds() {
	abilities := GetSeedsAbilities()
	configuratorAbilities := make(map[string]ConfiguratorAbility, len(abilities))
	for i, ability := range abilities {
		configuratorAbilities[i] = ConvertToConfiguratorAbility(*ability)
	}

	SaveAbilitiesToFile(configuratorAbilities)
	PlayersInitialAbilitiesIds := lo.Keys(configuratorAbilities)
	sort.Strings(PlayersInitialAbilitiesIds)
	SavePlayersInitialAbilities(PlayersInitialAbilitiesIds[1:5])
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
	if updatedAbility.ExecutionDurationMs != nil {
		ability.ExecutionDurationMs = updatedAbility.ExecutionDurationMs
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

	if _, exists := abilities[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ability not found"})
		return
	}

	if playersInitialAbilitiesIds, err := LoadPlayersInitialAbilitiesIds(); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "could not load players initial abilities"})
		return
	} else if lo.Contains(playersInitialAbilitiesIds, id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "ability is players initial ability, cannot be deleted, remove it as players initial ability to delete it"})
		return
	}

	delete(abilities, id)
	SaveAbilitiesToFile(abilities)
	c.JSON(http.StatusOK, gin.H{"message": "ability deleted"})
}

func GetPlayersInitialAbilities(c *gin.Context) {
	PlayersInitialAbilitiesIds, err := LoadPlayersInitialAbilitiesIds()
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, PlayersInitialAbilitiesIds)
}

func SetPlayersInitialAbilities(c *gin.Context) {
	var PlayersInitialAbilitiesIds []string
	if err := c.ShouldBindJSON(&PlayersInitialAbilitiesIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(PlayersInitialAbilitiesIds) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Selected more than 4 abilities"})
		return
	} else if len(PlayersInitialAbilitiesIds) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Less than 4 abilities where selected"})
		return
	}

	abilities, err := LoadAbilitiesFromFile()
	if err != nil {
		fmt.Printf("Error loading abilities: %v\n", err)
		return
	}

	notValidAbilitiesIds, _ := lo.Difference(PlayersInitialAbilitiesIds, lo.Keys(abilities))
	if len(notValidAbilitiesIds) > 0 {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Abilities IDs don't exist: %v", notValidAbilitiesIds)},
		)
		return
	}
	SavePlayersInitialAbilities(PlayersInitialAbilitiesIds)
	c.JSON(http.StatusOK, PlayersInitialAbilitiesIds)
}

func validateAbility(a ConfiguratorAbility) error {
	if a.Name == "" {
		return errors.New("`name` required")
	}
	if a.Cooldown == nil {
		return errors.New("`cooldown` required")
	}
	if a.ExecutionDurationMs == nil {
		return errors.New("`execution_duration_ms` required")
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
