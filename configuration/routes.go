package configuration

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer interface {
	gin.IRoutes
	Run(...string) error
}

func HealthCheck(context *gin.Context) {
	context.String(http.StatusOK, "OK, go!")
}

func Setup(server *gin.Engine) {
	server.HEAD("/health", HealthCheck)
}
