package configuration

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"os"
	"reflect"
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

		database := &gorm.DB{}
		monkey.Patch(gorm.Open, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
			return database, nil
		})

		expected := &sql.DB{}
		monkey.PatchInstanceMethod(reflect.TypeOf(database), "DB", func(*gorm.DB) (*sql.DB, error) {
			return expected, nil
		})

		// Act
		actual := ConnectToDatabase()

		// Assert
		assert.Equal(expected, actual)
	})

	test.Run("Should log a panic when failed to connect to database", func(test *testing.T) {
		// Arrange
		var capture bytes.Buffer
		log.SetOutput(&capture)
		monkey.Patch(gorm.Open, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
			return nil, errors.New("Failed to connect to database")
		})

		// Act
		actual := ConnectToDatabase()

		// Assert
		assert.Contains(capture.String(), "Failed to connect to database")
		assert.Nil(actual)
	})
}

func TestMigrateDatabase(test *testing.T) {
	// Database.AutoMigrate(&models.Book{})
}
