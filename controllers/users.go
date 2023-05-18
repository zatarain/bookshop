package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/bookshop/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
	Nickname string
	Password string
}

type UserController struct {
	Database *gorm.DB
}

func HashPassword(password string) (string, error) {
	hash, exception := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), exception
}

func (users *UserController) Signup(context *gin.Context) {
	var credentials Credentials

	// Trying to bind input from JSON
	if context.BindJSON(&credentials) != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read input",
		})
		return
	}

	// Trying to crete a hash for password
	hash, exception := HashPassword(credentials.Password)
	if exception != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to create the hash",
		})
		return
	}

	// Insert user into the database table users
	user := models.User{
		Nickname: credentials.Nickname,
		Password: hash,
	}
	users.Database.Create(&user)
	if user.ID != 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to insert user into table users",
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": "User successfully created",
	})
}
