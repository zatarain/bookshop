package configuration

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/zatarain/bookshop/mocks"
)

func TestSetup(test *testing.T) {
	test.Run("Should setup all the end-points", func(test *testing.T) {
		// Arrange
		server := new(mocks.MockedEngine)
		endPointHandler := mock.AnythingOfType("gin.HandlerFunc")
		autorisationHandler := mock.AnythingOfType("gin.HandlerFunc")
		server.On("HEAD", "/health", endPointHandler).Return(server)
		server.On("POST", "/signup", endPointHandler).Return(server)
		server.On("POST", "/login", endPointHandler).Return(server)
		server.On("GET", "/books", autorisationHandler, endPointHandler).Return(server)

		// Act
		Setup(server)

		// Assert
		server.AssertExpectations(test)
	})
}
