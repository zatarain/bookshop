package configuration

import (
	"log"
	"os"

	"github.com/zatarain/bookshop/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectToDatabase() {
	dialector := sqlite.Open(os.Getenv("DATABASE"))
	database, exception := gorm.Open(dialector, &gorm.Config{})
	if exception != nil {
		log.Panic("Failed to connect to the database.")
	}

	Database = database
}

func MigrateDatabase() {
	Database.AutoMigrate(&models.Book{})
}
