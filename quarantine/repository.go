package quarantine

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/quarantine/models"
)

type Repository interface {
	isExists(mobileNumber string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) isExists(mobileNumber string) (bool, error) {
	var user models.Quarantine
	err := r.db.Where(&models.Quarantine{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return false, nil
	}
	return user.MobileNumber == mobileNumber, nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}