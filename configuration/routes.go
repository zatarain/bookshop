package configuration

import (
	"github.com/gin-gonic/gin"
	"github.com/zatarain/bookshop/controllers"
)

func Setup(server gin.IRouter) {
	users := &controllers.UsersController{Database: Database}
	server.HEAD("/health", controllers.HealthCheck)
	server.POST("/signup", users.Signup)
	server.POST("/login", users.Login)
	server.GET("/books", users.Authorise, controllers.GetBooks)
}
