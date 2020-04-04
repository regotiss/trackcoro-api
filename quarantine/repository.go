package quarantine

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/quarantine/models"
)

type Repository interface {
	isExists(mobileNumber string) (bool, error)
	SaveDetails(quarantine models.Quarantine) error
}

type repository struct {
	db *gorm.DB
}

func (r repository) isExists(mobileNumber string) (bool, error) {
	user, err := r.getBy(mobileNumber)
	if err != nil {
		return false, err
	}
	return user.MobileNumber == mobileNumber, nil
}

func (r repository) SaveDetails(quarantine models.Quarantine) error {
	user, err := r.getBy(quarantine.MobileNumber)
	if err != nil || user.MobileNumber != quarantine.MobileNumber {
		logrus.Error("Quarantine does not exists ", err)
		return errors.New(NotExists)
	}
	quarantine.ID = user.ID
	err = r.db.Save(&quarantine).Error
	if err != nil {
		logrus.Error("Could not save details ", err)
	}
	return err
}

func (r repository) getBy(mobileNumber string) (models.Quarantine, error) {
	var user models.Quarantine
	err := r.db.Where(&models.Quarantine{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return models.Quarantine{}, err
	}
	return user, nil
}
func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
