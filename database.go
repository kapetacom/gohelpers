package gohelpers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() (*gorm.DB, error) {
	// Get environment variables
	host := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PW")
	schemaName := os.Getenv("POSTGRES_SCHEMA")

	// Create the initial connection string without specifying the schema
	dsn := fmt.Sprintf("host=%s user=postgres password=%s dbname=postgres sslmode=require", host, password)

	// Connect to the default 'postgres' database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check if the schema exists
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)", schemaName).Scan(&exists).Error
	if err != nil {
		return nil, fmt.Errorf("failed to check if schema exists: %w", err)
	}

	// If the schema doesn't exist, create it
	if !exists {
		err = db.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName)).Error
		if err != nil {
			return nil, fmt.Errorf("failed to create schema: %w", err)
		}
		log.Printf("Schema '%s' created successfully", schemaName)
	}

	// Close the initial connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	sqlDB.Close()

	// Connect to the database with the specific schema
	dsn = fmt.Sprintf("%s search_path=%s", dsn, schemaName)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with schema: %w", err)
	}

	log.Printf("Connected to database with schema '%s'", schemaName)
	return db, nil
}
