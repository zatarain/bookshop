package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zatarain/bookshop/mocks"
	"gorm.io/gorm"
)

func TestSignup(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

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

		// Act
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusCreated, recorder.Code)
		assert.Contains(recorder.Body.String(), "User successfully created")
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

		// Act
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "User already exists")
	})

}
