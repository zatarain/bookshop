package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zatarain/bookshop/mocks"
	"gorm.io/gorm"
)

func TestSignup(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	defer monkey.UnpatchAll()

	test.Run("Should create a new user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: nil})
		server.POST("/signup", users.Signup)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusCreated, recorder.Code)
		assert.Contains(recorder.Body.String(), "User successfully created")
		database.AssertExpectations(test)
	})

	test.Run("Should NOT create a duplicated user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: errors.New("User already exists")})
		server.POST("/signup", users.Signup)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "User already exists")
		database.AssertExpectations(test)
	})

	test.Run("Should NOT try to create a user when unable to bind JSON", func(test *testing.T) {
		// Arrange
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: nil})
		body := bytes.NewBuffer([]byte("Malformed JSON"))
		request, _ := http.NewRequest(http.MethodPost, "/signup", body)
		recorder := httptest.NewRecorder()
		server := gin.New()
		server.POST("/signup", users.Signup)

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to read input")
	})
}
