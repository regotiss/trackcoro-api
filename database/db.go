package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/database/models"
	"trackcoro/quarantine"
)

var DB *gorm.DB

func ConnectToDB() {
	logrus.Info("Connecting to DB")
	db, err := gorm.Open("postgres", os.Getenv(quarantine.DBConnectionString))
	if err != nil {
		logrus.Panic("Could not connect to db", err)
	}
	DB = db
	logrus.Info("DB connection established")
}

func MigrateSchema() {
	DB.AutoMigrate(&models.Quarantine{})
	DB.AutoMigrate(&models.QuarantineAddress{})
	DB.AutoMigrate(&models.QuarantineTravelHistory{})
}