package configuration

import (
	"github.com/gin-gonic/gin"
	"github.com/zatarain/bookshop/controllers"
)

func Setup(server gin.IRouter) {
	server.HEAD("/health", controllers.HealthCheck)
	server.GET("/books", controllers.GetBooks)
}