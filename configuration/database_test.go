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
	"github.com/zatarain/bookshop/models"
	"gorm.io/gorm"
)

func TestConnectToDatabase(test *testing.T) {
	assert := assert.New(test)
	monkey.Patch(log.Panic, log.Print)

	// Teardown test suite
	defer monkey.UnpatchAll()
	defer log.SetOutput(os.Stderr)

	test.Run("Should connect to database and return generic SQL connection pointer", func(test *testing.T) {
		// Arrange
		dummy := &gorm.DB{}
		monkey.Patch(gorm.Open, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
			return dummy, nil
		})

		expected := &sql.DB{}
		monkey.PatchInstanceMethod(reflect.TypeOf(dummy), "DB", func(*gorm.DB) (*sql.DB, error) {
			return expected, nil
		})

		// Act
		actual := ConnectToDatabase()

		// Assert
		assert.Equal(expected, actual)
		assert.Equal(dummy, Database)
	})

	test.Run("Should log a panic when failed to connect to database and return nil", func(test *testing.T) {
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

	test.Run("Should return nil when failed to get the generic SQL connection pointer", func(test *testing.T) {
		// Arrange
		var capture bytes.Buffer
		log.SetOutput(&capture)

		database := &gorm.DB{}
		monkey.Patch(gorm.Open, func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error) {
			return database, nil
		})

		monkey.PatchInstanceMethod(reflect.TypeOf(database), "DB", func(*gorm.DB) (*sql.DB, error) {
			return nil, errors.New("Failed to get SQL connection pointer")
		})

		// Act
		actual := ConnectToDatabase()

		// Assert
		assert.Contains(capture.String(), "Failed to get SQL connection pointer")
		assert.Nil(actual)
	})
}

func TestMigrateDatabase(test *testing.T) {
	assert := assert.New(test)
	monkey.Patch(log.Panic, log.Print)

	// Teardown test suite
	defer monkey.UnpatchAll()
	defer log.SetOutput(os.Stderr)

	test.Run("Should connect to database and return generic SQL connection pointer", func(test *testing.T) {
		// Arrange
		ConnectToDatabase()

		// Act
		MigrateDatabase()

		// Assert
		Database.Create(&models.Book{
			Title:    "dummy book",
			Author:   "dummy author",
			Price:    100,
			Quantity: 4,
		})
		var book models.Book
		Database.Find(&book, 1)
		assert.Equal(book.Title, "dummy book")
	})
}
