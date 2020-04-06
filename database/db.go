package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/constants"
	"trackcoro/database/models"
)

var DB *gorm.DB

func ConnectToDB() {
	logrus.Info("Connecting to DB")
	db, err := gorm.Open("postgres", os.Getenv(constants.DBConnectionString))
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

	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.SupervisingOfficer{})

	DB.Model(&models.SupervisingOfficer{}).AddForeignKey("admin_id", "admins(id)", "CASCADE", "NO ACTION")
	DB.Model(&models.Quarantine{}).AddForeignKey("supervising_officer_id", "supervising_officers(id)", "CASCADE", "NO ACTION")
}