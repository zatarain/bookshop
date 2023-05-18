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
	"github.com/zatarain/bookshop/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestSignup(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	// Teardown test suite
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
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: nil})
		server.POST("/signup", users.Signup)
		body := bytes.NewBuffer([]byte("Malformed JSON"))
		request, _ := http.NewRequest(http.MethodPost, "/signup", body)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to read input")
		database.AssertNotCalled(test, "Create", mock.AnythingOfType("*models.User"))
	})

	test.Run("Should NOT try to create a user when unable hash password", func(test *testing.T) {
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
		monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
			return []byte{}, errors.New("Unable to hash")
		})

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to create the hash for password")
		database.AssertNotCalled(test, "Create", mock.AnythingOfType("*models.User"))
	})
}

func TestLogin(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should login the user and create the token", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		call := database.On(
			"First",
			mock.AnythingOfType("*models.User"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			user := arguments.Get(0).(*models.User)
			user.ID = 12345
			user.Nickname = "dummy-user"
			user.Password = "top-secret"
		}

		calledToCompareHashAndPassword := false
		monkey.Patch(bcrypt.CompareHashAndPassword, func([]byte, []byte) error {
			calledToCompareHashAndPassword = true
			return nil
		})

		server.POST("/login", users.Login)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.True(calledToCompareHashAndPassword)
		//assert.Equal(http.StatusOK, recorder.Code)
		assert.Contains(recorder.Body.String(), "Yaaay! You are logged in :)")
		database.AssertExpectations(test)
	})
}
