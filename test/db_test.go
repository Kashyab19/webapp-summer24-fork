// db_test.go
package test

import (
	"os"
	"testing"
	"webapp/models" // Adjust this import path to your models package

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDatabaseSetup tests the database connection, migration, and optional seeding.
func TestDatabaseSetup(t *testing.T) {
	// Retrieve the database connection string from environment variables or use a fallback for testing
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=kashyabmurali dbname=postgres password=postgres port=5432 sslmode=disable"
	}

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// if err := db.Debug().AutoMigrate(&models.User{}); err != nil {
	// 	t.Fatalf("Failed to auto-migrate: %v", err)
	// }

	seedDatabase(db, t)

	// Close the database connection at the end of the test
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get database connection: %v", err)
	}
	defer sqlDB.Close()
}

// seedDatabase is an optional helper function to seed the database with initial data.
func seedDatabase(db *gorm.DB, t *testing.T) {
	// Example of seeding a user - customize according to your model definitions
	user := models.User{
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "securepassword", // Note: In a real application, ensure passwords are securely hashed
	}

	// Create the user in the database
	if result := db.Create(&user); result.Error != nil {
		t.Logf("Failed to seed database with user: %v", result.Error)
	}
}
