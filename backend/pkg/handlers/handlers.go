package handlers

import (
	"backend/pkg/game/repository"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() {
	r := gin.Default()

	r.GET("/characters", getCharactersHandler)

	const HTTP_PORT string = "8081"
	go func() {
		if os.Getenv("GO_ENV") != "test" {
			fmt.Println("HTTP server running on port", HTTP_PORT)
		}
		if err := r.Run(":" + HTTP_PORT); err != nil {
			panic(err)
		}
	}()
}

func getCharactersHandler(c *gin.Context) {
	chars, err := repository.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch characters"})
		return
	}
	c.JSON(http.StatusOK, chars)
}
