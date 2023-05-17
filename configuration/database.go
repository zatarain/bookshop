package configuration

import (
	"log"

	"github.com/zatarain/bookshop/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() {
	database, exception := gorm.Open(sqlite.Open("data/test.db"), &gorm.Config{})
	if exception != nil {
		log.Panic("Failed to connect to the database.")
	}

	Database = database
}

func Migrate() {
	Database.AutoMigrate(&models.Book{})
}
