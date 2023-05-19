package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zatarain/bookshop/mocks"
	"github.com/zatarain/bookshop/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
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

	CompareSuccessful := func([]byte, []byte) error {
		return nil
	}

	CompareFailure := func([]byte, []byte) error {
		return errors.New("Invalid Password")
	}

	NiceFakeToken := func(*UsersController, *gin.Context, *models.User) (string, error) {
		return "Nice Fake Token", nil
	}

	NoToken := func(*UsersController, *gin.Context, *models.User) (string, error) {
		return "", errors.New("No Token")
	}

	CheckCookie := func(cookie *http.Cookie) bool {
		log.Printf("Cookie => %v\n", *cookie)
		return cookie.Name == "Authorisation"
	}

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should login the user and create the token", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		anyUser := mock.AnythingOfType("*models.User")
		call := database.
			On("First", anyUser, "nickname = ?", "dummy-user").
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			user := arguments.Get(0).(*models.User)
			user.ID = 12345
			user.Nickname = "dummy-user"
			user.Password = "top-secret"
		}

		monkey.Patch(bcrypt.CompareHashAndPassword, CompareSuccessful)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "NewToken", NiceFakeToken)
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
		cookies := recorder.Result().Cookies()
		index := slices.IndexFunc(cookies, CheckCookie)

		// Assert
		database.AssertExpectations(test)
		assert.Equal(http.StatusOK, recorder.Code)
		assert.Contains(recorder.Body.String(), "Yaaay! You are logged in :)")
		require.GreaterOrEqual(test, index, 0)
		assert.Equal("Nice+Fake+Token", cookies[index].Value)
		assert.Equal(7*24*60*60, cookies[index].MaxAge)
		assert.False(cookies[index].Secure)
		assert.True(cookies[index].HttpOnly)
	})

	test.Run("Should response with internal server error when unable to generate token", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		anyUser := mock.AnythingOfType("*models.User")
		call := database.
			On("First", anyUser, "nickname = ?", "dummy-user").
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			user := arguments.Get(0).(*models.User)
			user.ID = 12345
			user.Nickname = "dummy-user"
			user.Password = "top-secret"
		}

		monkey.Patch(bcrypt.CompareHashAndPassword, CompareSuccessful)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "NewToken", NoToken)
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
		cookies := recorder.Result().Cookies()
		index := slices.IndexFunc(cookies, CheckCookie)

		// Assert
		database.AssertExpectations(test)
		assert.Equal(http.StatusInternalServerError, recorder.Code)
		assert.Contains(recorder.Body.String(), "Unable to generate access token")
		require.Equal(test, index, -1)
	})

	test.Run("Should NOT try to login the user when unable to bind JSON", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		anyUser := mock.AnythingOfType("*models.User")
		database.On("First", anyUser).Return(&gorm.DB{Error: nil})
		server.POST("/login", users.Login)
		body := bytes.NewBuffer([]byte("Malformed JSON"))
		request, _ := http.NewRequest(http.MethodPost, "/login", body)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to read input")
		database.AssertNotCalled(test, "First", mock.AnythingOfType("*models.User"))
	})

	user := Credentials{
		Nickname: "dummy-user",
		Password: "top-secret",
	}

	InvalidNicknameOrPasswordTestcases := []struct {
		description string
		user        models.User
		compare     func([]byte, []byte) error
	}{
		{
			description: "Should NOT login the user when we didn't find nickname in database",
			user: models.User{
				ID:       0,
				Nickname: "",
				Password: "",
			},
			compare: CompareSuccessful,
		},
		{
			description: "Should NOT login the user when password doesn't match with stored hash",
			user: models.User{
				ID:       12345,
				Nickname: user.Nickname,
				Password: "secret-top",
			},
			compare: CompareFailure,
		},
		{
			description: "Should NOT login the user when either we didn't find nickname in database or password doesn't match",
			user: models.User{
				ID:       0,
				Nickname: "",
				Password: "",
			},
			compare: CompareFailure,
		},
	}

	anyUser := mock.AnythingOfType("*models.User")

	for _, testcase := range InvalidNicknameOrPasswordTestcases {
		test.Run(testcase.description, func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			users := &UsersController{Database: database}
			call := database.
				On("First", anyUser, "nickname = ?", user.Nickname).
				Return(&gorm.DB{Error: nil})
			call.RunFn = func(arguments mock.Arguments) {
				user := arguments.Get(0).(*models.User)
				user.ID = testcase.user.ID
				user.Nickname = testcase.user.Nickname
				user.Password = testcase.user.Password
			}

			monkey.Patch(bcrypt.CompareHashAndPassword, testcase.compare)

			server.POST("/login", users.Login)
			body, _ := json.Marshal(user)
			request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusBadRequest, recorder.Code)
			assert.Contains(recorder.Body.String(), "Invalid nickname or password")
			database.AssertExpectations(test)
		})
	}
}

func TestAuthorise(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	today := time.Now()
	dummy := models.User{
		ID:        12345,
		Nickname:  "dummy-user",
		Password:  "top-secret",
		CreatedAt: today,
		UpdatedAt: today,
	}

	AuthorisedEndPointHandler := func(context *gin.Context) {
		data, exists := context.Get("user")
		user := data.(*models.User)

		assert.True(exists)
		assert.Equal(&dummy, user)
	}

	UnauthorisedEndPointHandler := func(*gin.Context) {
		assert.False(true, "This should never run otherwise the test failed!")
	}

	ValidToken := func(*UsersController, *gin.Context) (*models.User, error) {
		return &dummy, nil
	}

	InvalidToken := func(*UsersController, *gin.Context) (*models.User, error) {
		return nil, errors.New("Invalid token")
	}

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should set the user within the context and continue when token is valid", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		server.GET("/", users.Authorise, AuthorisedEndPointHandler)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "ValidateToken", ValidToken)
		request, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusOK, recorder.Code)
	})

	test.Run("Should set the user within the context and continue when token is valid", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		server.GET("/", users.Authorise, UnauthorisedEndPointHandler)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "ValidateToken", InvalidToken)
		request, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusUnauthorized, recorder.Code)
		assert.Contains(recorder.Body.String(), "Invalid token")
	})
}
