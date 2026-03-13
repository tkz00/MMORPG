package repository

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectPostgres() error {
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

	if user == "" || password == "" || dbname == "" {
		log.Println("[WARN] One or more database env vars are empty.")
		log.Printf("POSTGRES_USER=%q POSTGRES_DB=%q POSTGRES_PORT=%q\n", user, dbname, port)
	}

	// fail loudly if someone left defaults
	if password == "changeme" {
		return fmt.Errorf("POSTGRES_PASSWORD is still 'changeme', please configure .env")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	if err := waitForPostgres(dsn, 10*time.Second); err != nil {
		panic("postgres never became ready")
	}

	db, err := gorm.Open(postgres.Open(dsn), NewGormConfig())
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Auto-create tables for DB models
	if err := db.AutoMigrate(
		&CharacterDB{},
		&EffectDB{},
		&CharacterPerkDB{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	DB = db
	return nil
}

func waitForPostgres(dsn string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			db.Close()
			if err == nil {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for postgres")
}

func NewGormConfig() *gorm.Config {
	var logLevel logger.LogLevel
	if os.Getenv("GO_ENV") == "test" {
		logLevel = logger.Silent
		log.SetOutput(io.Discard)
	} else {
		logLevel = logger.Warn
	}

	return &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}
}
