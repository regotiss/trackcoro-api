package admin

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	Add(admin models.Admin) error
}

type repository struct {
	db *gorm.DB
}

func (r repository) Add(admin models.Admin) error {
	existingAdmin, err := r.getBy(admin.MobileNumber)
	if err != nil {
		logrus.Info("Adding admin")
		r.db.Create(&admin)
		return err
	}
	logrus.Info("Admin already exists")
	admin.ID = existingAdmin.ID
	err = r.db.Save(admin).Error
	return nil
}

func (r repository) IsExists(mobileNumber string) bool {
	user, err := r.getBy(mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) getBy(mobileNumber string) (models.Admin, error) {
	var user models.Admin
	err := r.db.Where(&models.Admin{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return models.Admin{}, errors.New(constants.NotExists)
	}
	return user, nil
}


func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
