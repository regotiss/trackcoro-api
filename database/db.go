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
	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.SupervisingOfficer{})
	DB.Model(&models.SupervisingOfficer{}).AddForeignKey("admin_id", "admins(id)", "CASCADE", "NO ACTION")

	DB.AutoMigrate(&models.Quarantine{})
	DB.Model(&models.Quarantine{}).AddForeignKey("supervising_officer_id", "supervising_officers(id)", "RESTRICT", "NO ACTION")

	DB.AutoMigrate(&models.QuarantineAddress{})
	DB.Model(&models.QuarantineAddress{}).AddForeignKey("quarantine_id", "quarantines(id)", "CASCADE", "NO ACTION")

	DB.AutoMigrate(&models.QuarantineTravelHistory{})
	DB.Model(&models.QuarantineTravelHistory{}).AddForeignKey("quarantine_id", "quarantines(id)", "CASCADE", "NO ACTION")

	DB.AutoMigrate(&models.PhotoUpload{})
	DB.Model(&models.PhotoUpload{}).AddForeignKey("quarantine_id", "quarantines(id)", "CASCADE", "NO ACTION")
}