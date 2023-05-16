package configuration

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestLoad(test *testing.T) {
	// Arrange
	assert := assert.New(test)
	isLoaded := false
	monkey.Patch(godotenv.Load, func(...string) error {
		isLoaded = true
		return nil
	})

	// Act
	Load()

	// Assert
	assert.True(isLoaded)

	// Teardown
	monkey.UnpatchAll()
}

func TestLoad_Error(test *testing.T) {
	// Arrange
	assert := assert.New(test)
	hasBeenCalled := false
	monkey.Patch(godotenv.Load, func(...string) error {
		hasBeenCalled = true
		return errors.New("Error while loading")
	})

	monkey.Patch(log.Fatal, log.Print)

	var capture bytes.Buffer
	log.SetOutput(&capture)

	// Act
	Load()

	// Assert
	assert.Equal(true, hasBeenCalled)
	assert.Contains(capture.String(), "Error loading environment variables file.")

	// Teardown
	log.SetOutput(os.Stderr)
	monkey.UnpatchAll()
}
