package configuration

import (
	"database/sql"
	"log"
	"os"

	"github.com/zatarain/bookshop/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectToDatabase() *sql.DB {
	filename := os.ExpandEnv("$GOPATH/src/$MODNAME/$DATABASE")
	dialector := sqlite.Open(filename)
	database, exception := gorm.Open(dialector, &gorm.Config{})
	if exception != nil {
		log.Panic("Failed to connect to the database.", exception.Error())
		return nil
	}

	connection, exception := database.DB()
	if exception != nil {
		log.Panic("Failed to get generic SQL connection pointer.", exception.Error())
		return nil
	}

	Database = database
	return connection
}

func MigrateDatabase() {
	Database.AutoMigrate(&models.Book{})
}
