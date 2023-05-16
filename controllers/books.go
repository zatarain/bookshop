package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "hello"})
}
