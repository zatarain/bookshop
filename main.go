package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zatarain/bookshop/configuration"
)

func main() {
	// Load Environment Variables
	configuration.Load()

	// Connect to Database
	connection := configuration.ConnectToDatabase()
	defer connection.Close()

	// Initialise Database
	configuration.MigrateDatabase()

	// Initialise the API Server
	server := gin.Default()
	configuration.Setup(server)
	if exception := server.Run(); exception != nil {
		log.Panic(exception.Error())
	}
}
