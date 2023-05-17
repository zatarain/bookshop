package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zatarain/bookshop/configuration"
)

func TestMain(test *testing.T) {
	ENVIRONMENT := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	monkey.Patch(log.Panic, log.Print)

	// Teardown test suite
	defer monkey.UnpatchAll()
	defer log.SetOutput(os.Stderr)
	defer os.Setenv("ENVIRONMENT", ENVIRONMENT)

	test.Run("Should run the service", func(test *testing.T) {
		// Arrange
		environmentHasBeenLoaded := false
		serverHasBeenSetup := false
		serverIsRunning := false
		monkey.Patch(configuration.Load, func() {
			environmentHasBeenLoaded = true
		})
		monkey.Patch(configuration.Setup, func(server gin.IRouter) {
			serverHasBeenSetup = true
			monkey.PatchInstanceMethod(reflect.TypeOf(server), "Run", func(*gin.Engine, ...string) error {
				serverIsRunning = true
				return nil
			})
		})

		// Act
		main()

		// Assert
		assert.True(environmentHasBeenLoaded)
		assert.True(serverHasBeenSetup)
		assert.True(serverIsRunning)
	})

	test.Run("Should log panic when failed to run server", func(test *testing.T) {
		// Arrange
		var capture bytes.Buffer
		log.SetOutput(&capture)
		monkey.Patch(configuration.Setup, func(server gin.IRouter) {
			monkey.PatchInstanceMethod(reflect.TypeOf(server), "Run", func(*gin.Engine, ...string) error {
				return errors.New("Failed to start the server.")
			})
		})

		// Act
		main()

		// Assert
		assert.Contains(capture.String(), "Failed to start the server.")
	})
}
