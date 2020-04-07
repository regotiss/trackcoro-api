package utils

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
)

func GetSOBy(db *gorm.DB, mobileNumber string) (models.SupervisingOfficer, error) {
	var user models.SupervisingOfficer
	err := db.Where(&models.SupervisingOfficer{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not find mobile number in db ", err)
		return models.SupervisingOfficer{}, errors.New(constants.QuarantineNotExistsError)
	}
	return user, nil
}

func GetQuarantineBy(db *gorm.DB, mobileNumber string) (models.Quarantine, error) {
	var user models.Quarantine
	err := db.Where(&models.Quarantine{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return models.Quarantine{}, errors.New(constants.QuarantineNotExistsError)
	}
	return user, nil
}

func GetQuarantines(db *gorm.DB, soMobileNumber string) ([]models.Quarantine, error) {
	existingSO, err := GetSOBy(db, soMobileNumber)
	if err != nil {
		return nil, err
	}
	var Quarantines []models.Quarantine
	err = db.Model(&existingSO).Related(&Quarantines).Error
	return Quarantines, err
}
