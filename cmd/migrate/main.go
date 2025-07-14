package main

import (
	"flag"
	"fmt"
	"log"

	"meteo-api/internal/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	var direction = flag.String("direction", "up", "Migration direction (up/down)")
	var steps = flag.Int("steps", 0, "Number of migration steps (0 for all)")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	dbConfig := database.NewConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migration instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Run migrations
	switch *direction {
	case "up":
		if *steps == 0 {
			err = m.Up()
		} else {
			err = m.Steps(*steps)
		}
	case "down":
		if *steps == 0 {
			err = m.Down()
		} else {
			err = m.Steps(-*steps)
		}
	default:
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", *direction)
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No migrations to apply")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Printf("Migration %s completed successfully\n", *direction)
	}
}