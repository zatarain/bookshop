package configuration

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestConnectToDatabase(test *testing.T) {
	ENVIRONMENT := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	assert := assert.New(test)
	monkey.Patch(log.Panic, log.Print)

	// Teardown test suite
	defer monkey.UnpatchAll()
	defer log.SetOutput(os.Stderr)
	defer os.Setenv("ENVIRONMENT", ENVIRONMENT)

	test.Run("Should log a panic when failed to connect to database", func(test *testing.T) {
		// Arrange
		var capture bytes.Buffer
		log.SetOutput(&capture)
		monkey.Patch(gorm.Open, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
			return nil, errors.New("Failed to connect to database.")
		})

		// Act
		ConnectToDatabase()

		// Assert
		assert.Contains(capture.String(), "Failed to connect to database.")
	})
}

func TestMigrateDatabase(test *testing.T) {
	// Database.AutoMigrate(&models.Book{})
}
