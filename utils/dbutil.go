package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
	models2 "trackcoro/models"
)

func GetSOBy(db *gorm.DB, mobileNumber string) (models.SupervisingOfficer, *models2.Error) {
	var user models.SupervisingOfficer
	err := db.Where(&models.SupervisingOfficer{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("SO not found with given mobile number ", err)
		return models.SupervisingOfficer{}, &constants.SONotExistsError
	}
	return user, nil
}

func GetAllQuarantineDetails(db *gorm.DB, mobileNumber string) (models.Quarantine, error) {
	var user models.Quarantine
	err := db.Preload("SupervisingOfficer").Preload("Address").Preload("TravelHistory").Where(&models.Quarantine{MobileNumber: mobileNumber}).
		First(&user).Error
	if err != nil {
		logrus.Error("Quarantine not found with given mobile number ", err)
		return models.Quarantine{}, constants.QuarantineNotExistsError
	}
	return user, nil
}

func GetSOByID(db *gorm.DB, ID uint) (models.SupervisingOfficer, error) {
	var user models.SupervisingOfficer
	err := db.Where(&models.SupervisingOfficer{Model: gorm.Model{ID: ID}}).First(&user).Error
	if err!=nil{
		logrus.Error("Could not find supervisor with given ID", err)
		return models.SupervisingOfficer{}, constants.SONotExistsError
	}
	return user, nil
}

func GetQuarantineBy(db *gorm.DB, mobileNumber string) (models.Quarantine, error) {
	var user models.Quarantine
	err := db.Where(&models.Quarantine{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Quarantine not found with given mobile number ", err)
		return models.Quarantine{}, constants.QuarantineNotExistsError
	}
	return user, nil
}

func GetQuarantines(db *gorm.DB, soMobileNumber string) ([]models.Quarantine, *models2.Error) {
	existingSO, err := GetSOBy(db, soMobileNumber)
	if err != nil {
		return nil, err
	}
	var Quarantines []models.Quarantine
	dbError := db.Model(&existingSO).Preload("Address").Preload("TravelHistory").Related(&Quarantines).Error
	if dbError != nil {
		return nil, &constants.InternalError
	}
	return Quarantines, nil
}
