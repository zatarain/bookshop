package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zatarain/bookshop/models"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Nickname string
	Password string
}

type UsersController struct {
	Database models.DataAccessInterface
}

func (credentials *Credentials) HashPassword() error {
	hash, exception := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	credentials.Password = string(hash)
	return exception
}

func getCredentialsFromRequest(context *gin.Context) *Credentials {
	var credentials Credentials

	// Trying to bind input from JSON
	if binding := context.BindJSON(&credentials); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to read input",
			"details": binding.Error(),
		})
		return nil
	}

	return &credentials
}

func (users *UsersController) Signup(context *gin.Context) {
	credentials := getCredentialsFromRequest(context)
	if credentials == nil {
		return
	}

	// Trying to crete a hash for password
	if exception := credentials.HashPassword(); exception != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to create the hash for password",
			"details": exception.Error(),
		})
		return
	}

	// Insert user into the database table users
	user := models.User{
		Nickname: credentials.Nickname,
		Password: credentials.Password,
	}
	inserting := users.Database.Create(&user).Error
	if inserting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to insert user into table users",
			"details": inserting.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"summary": "User successfully created",
		"details": user.String(),
	})
}

func generateTokenAndCookie(context *gin.Context, nickname string) {
	// Create the token
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": nickname,
			"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	)

	// Signing the token with secret key
	signedToken, exception := token.SignedString(os.Getenv("SECRET_TOKEN_KEY"))
	if exception != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"summary": "Unable to generate access token",
			"details": exception.Error(),
		})
	}

	// Send cookie to the client
	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorisation", signedToken, 7*24*60*60, "", "", false, true)
}

func (users *UsersController) Login(context *gin.Context) {
	credentials := getCredentialsFromRequest(context)
	if credentials == nil {
		return
	}

	// Checking the credentials
	var user models.User
	users.Database.First(&user, "nickname = ?", credentials.Nickname)
	failed := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(credentials.Password),
	)
	if user.ID == 0 || failed != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Invalid nickname or password",
		})
		return
	}

	// Generate JWT Token and send it in the Cookies
	generateTokenAndCookie(context, user.Nickname)
	context.JSON(http.StatusOK, gin.H{"summary": "Yaaay! You are logged In :)"})
}

func (users *UsersController) Authorise(context *gin.Context) {
	context.Next()
}
