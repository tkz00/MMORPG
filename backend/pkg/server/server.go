package server

import (
	"backend/connection"
	"backend/pkg/configurator"
	"backend/pkg/game/repository"
	"backend/pkg/handlers"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type ServerOptions struct {
	EnvFilePath string
	Quiet       bool
}

func Start(opts ServerOptions) error {
	// Load .env if given
	if opts.EnvFilePath != "" {
		if err := godotenv.Load(opts.EnvFilePath); err != nil {
			fmt.Println("Warning: failed to load env:", err)
		}
	}

	// Silence logs in test mode
	if opts.Quiet {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}

	// Run configurator seeding and server
	configurator.RunSeeds()
	go configurator.Run()

	// Connect DB
	if err := repository.ConnectPostgres(); err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	fmt.Println("Postgres connected successfully")

	repository.RunSeeds()
	handlers.RegisterRoutes()

	// Start websocket server
	const PORT = "3009"
	go connection.NewServer().StartConnection(PORT)

	select {} // block forever
}
