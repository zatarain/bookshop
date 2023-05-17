package configuration

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	filename := fmt.Sprintf("%s.env", os.Getenv("ENVIRONMENT"))
	log.Println("Loading environment file:", filename)
	exception := godotenv.Load(filename)
	if exception != nil {
		log.Panic("Error loading environment variables file.", exception.Error())
	}
}
