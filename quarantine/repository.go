package quarantine

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
	"trackcoro/quarantine/models"
)

type Repository interface {
	isExists(mobileNumber string) bool
	SaveDetails(quarantine models.Quarantine) error
	GetQuarantineDays(mobileNumber string) (uint, time.Time, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) isExists(mobileNumber string) bool {
	user, err := r.getBy(mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) SaveDetails(quarantine models.Quarantine) error {
	user, err := r.getBy(quarantine.MobileNumber)
	if err != nil {
		return err
	}
	quarantine.ID = user.ID
	err = r.db.Save(&quarantine).Error
	if err != nil {
		logrus.Error("Could not save details ", err)
	}
	return err
}

func (r repository) GetQuarantineDays(mobileNumber string) (uint, time.Time, error) {
	user, err := r.getBy(mobileNumber)
	if err != nil {
		return 0, time.Time{}, err
	}
	return user.NoOfQuarantineDays, user.QuarantineStartedFrom, nil
}


func (r repository) getBy(mobileNumber string) (models.Quarantine, error) {
	var user models.Quarantine
	err := r.db.Where(&models.Quarantine{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return models.Quarantine{}, errors.New(NotExists)
	}
	return user, nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
