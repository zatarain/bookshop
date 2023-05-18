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
		server.On("HEAD", "/health", mock.AnythingOfType("gin.HandlerFunc")).Return(server)
		server.On("GET", "/books", mock.AnythingOfType("gin.HandlerFunc")).Return(server)
		server.On("POST", "/signup", mock.AnythingOfType("gin.HandlerFunc")).Return(server)

		// Act
		Setup(server)

		// Assert
		server.AssertExpectations(test)
	})
}
