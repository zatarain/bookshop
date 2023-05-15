package main

import (
	"bookshop/configuration"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("OK, go!")
	defer log.Println("Good bye!")
	configuration.Load()
	server := gin.Default()
	configuration.Setup(server)
	server.Run()
}
