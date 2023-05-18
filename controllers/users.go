package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func HashPassword(password string) (string, error) {
	hash, exception := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), exception
}

func (users *UsersController) Signup(context *gin.Context) {
	var credentials Credentials

	// Trying to bind input from JSON
	binding := context.BindJSON(&credentials)
	if binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to read input",
			"details": binding.Error(),
		})
		return
	}

	// Trying to crete a hash for password
	hash, exception := HashPassword(credentials.Password)
	if exception != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to create the hash",
			"details": exception.Error(),
		})
		return
	}

	// Insert user into the database table users
	user := models.User{
		Nickname: credentials.Nickname,
		Password: hash,
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
