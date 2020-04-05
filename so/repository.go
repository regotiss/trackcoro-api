package so

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	AddQuarantine(mobileNumber string, quarantine models.Quarantine) error
	GetQuarantines(mobileNumber string) ([]models.Quarantine, error)
}

type repository struct {
	db *gorm.DB
}


func (r repository) IsExists(mobileNumber string) bool {
	user, err := r.getBy(mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) AddQuarantine(mobileNumber string, quarantine models.Quarantine) error {
	existingSO, err := r.getBy(mobileNumber)
	if err != nil {
		return err
	}
	quarantine.SupervisingOfficerID = existingSO.ID
	return r.db.Save(&quarantine).Error
}

func (r repository) GetQuarantines(mobileNumber string) ([]models.Quarantine, error) {
	existingSO, err := r.getBy(mobileNumber)
	if err != nil {
		return nil, err
	}
	var Quarantines []models.Quarantine
	err = r.db.Model(&existingSO).Related(&Quarantines).Error
	return Quarantines, err
}


func (r repository) getBy(mobileNumber string) (models.SupervisingOfficer, error) {
	var user models.SupervisingOfficer
	err := r.db.Where(&models.SupervisingOfficer{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not find mobile number in db ", err)
		return models.SupervisingOfficer{}, errors.New(constants.NotExists)
	}
	return user, nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
