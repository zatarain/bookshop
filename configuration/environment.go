package configuration

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() {
	exception := godotenv.Load()
	if exception != nil {
		log.Fatal("Error loading environment variables file.")
	}
}
