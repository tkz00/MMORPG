package repository

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() error {
	// Load variables from .env if present
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println(err)
	}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	// fail loudly if someone left defaults
	if password == "changeme" {
		return fmt.Errorf("POSTGRES_PASSWORD is still 'changeme', please configure .env")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Auto-create tables for DB models
	if err := db.AutoMigrate(&CharacterDB{}); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	return nil
}
