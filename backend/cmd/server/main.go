// Copyright (C) 2025 Theo Katz
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"backend/connection"
	"backend/pkg/configurator"
	"backend/pkg/game/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	configurator.RunSeeds()

	// Start configurator server
	go configurator.Run()

	if err := repository.ConnectPostgres(); err != nil {
		panic(err)
	}
	fmt.Println("Postgres connected successfully")

	repository.RunSeeds()

	SetupCharactersRouter()

	// Start game server
	const PORT string = "3009"
	server := connection.CreateServer()
	go server.StartConnection(PORT)

	// Block main from exiting by waiting indefinitely
	select {}
}

func SetupCharactersRouter() {
	r := gin.Default()

	r.GET("/characters", func(c *gin.Context) {
		chars, err := repository.GetAllCharacters() // implement this in your repo
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch characters"})
			return
		}
		c.JSON(http.StatusOK, chars)
	})

	const HTTP_PORT string = "8081"
	go func() {
		fmt.Println("HTTP server running on port", HTTP_PORT)
		if err := r.Run(":" + HTTP_PORT); err != nil {
			panic(err)
		}
	}()
}
