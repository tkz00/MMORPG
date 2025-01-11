package configurator

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// Register your routes here
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return r
}

func Run() {
	r := SetupRouter()
	fmt.Println("Starting server on port 8080...")
	r.Run() // Default port
}
