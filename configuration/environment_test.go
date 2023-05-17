package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestLoad(test *testing.T) {
	assert := assert.New(test)
	test.Run("Should load the environment with godotenv", func(test *testing.T) {
		// Arrange
		ENVIRONMENT := os.Getenv("ENVIRONMENT")
		os.Setenv("ENVIRONMENT", "test")
		isLoaded := false
		var filenames []string

		monkey.Patch(godotenv.Load, func(envfiles ...string) error {
			filenames = slices.Clone(envfiles)
			isLoaded = true
			return nil
		})

		// Act
		Load()

		// Assert
		fmt.Println("filenames in Test:", filenames)
		//assert.EqualValues([]string{"test.env"}, filenames)
		assert.True(isLoaded)

		// Teardown
		os.Setenv("ENVIRONMENT", ENVIRONMENT)
		monkey.UnpatchAll()
	})

	test.Run("Should log the error as panic while trying to load the environment", func(test *testing.T) {
		// Arrange
		hasBeenCalled := false
		monkey.Patch(log.Panic, log.Print)
		var capture bytes.Buffer
		log.SetOutput(&capture)
		monkey.Patch(godotenv.Load, func(...string) error {
			hasBeenCalled = true
			return errors.New("Error while loading")
		})

		// Act
		Load()

		// Assert
		assert.Equal(true, hasBeenCalled)
		assert.Contains(capture.String(), "Error loading environment variables file.")

		// Teardown
		log.SetOutput(os.Stderr)
		monkey.UnpatchAll()
	})
}
