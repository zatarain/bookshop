package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zatarain/bookshop/configuration"
)

func main() {
	configuration.Load()
	configuration.ConnectToDatabase()
	configuration.MigrateDatabase()
	server := gin.Default()
	configuration.Setup(server)
	server.Run()
}
