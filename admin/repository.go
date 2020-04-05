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
	AddSO(adminMobileNumber string, so models.SupervisingOfficer) error
	GetSOs(adminMobileNumber string) ([]models.SupervisingOfficer, error)
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

func (r repository) Add(admin models.Admin) error {
	existingAdmin, err := r.getBy(admin.MobileNumber)
	if err != nil {
		logrus.Info("Adding admin")
		return r.db.Create(&admin).Error
	}
	logrus.Info("Admin already exists")
	admin.ID = existingAdmin.ID
	return r.db.Save(&admin).Error
}

func (r repository) AddSO(mobileNumber string, so models.SupervisingOfficer) error {
	existingAdmin, err := r.getBy(mobileNumber)
	if err != nil {
		return err
	}
	so.AdminID = existingAdmin.ID
	return r.db.Save(&so).Error
}

func (r repository) GetSOs(adminMobileNumber string) ([]models.SupervisingOfficer, error) {
	existingAdmin, err := r.getBy(adminMobileNumber)
	if err != nil {
		return nil, err
	}
	var SOs []models.SupervisingOfficer
	err = r.db.Model(&existingAdmin).Related(&SOs).Error
	return SOs, err
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
