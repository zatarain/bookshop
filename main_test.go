package main

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/zatarain/bookshop/configuration"
)

func TestMain(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	defer monkey.UnpatchAll()

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
}
