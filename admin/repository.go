package admin

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
	"trackcoro/utils"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	Add(admin models.Admin) error
	AddSO(adminMobileNumber string, so models.SupervisingOfficer) error
	GetSOs(adminMobileNumber string) ([]models.SupervisingOfficer, error)
	GetQuarantines(adminMobileNumber string, soMobileNumber string) ([]models.Quarantine, error)
	DeleteSO(adminMobileNumber string, soMobileNumber string) error
	ReplaceSO(adminMobileNumber string, oldSOMobileNumber string, newSOMobileNumber string) error
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

func (r repository) GetQuarantines(adminMobileNumber string, soMobileNumber string) ([]models.Quarantine, error) {
	_, _, err := r.isAdminOfSO(adminMobileNumber, soMobileNumber)
	if err != nil {
		return nil, err
	}
	return utils.GetQuarantines(r.db, soMobileNumber)
}

func (r repository) DeleteSO(adminMobileNumber string, soMobileNumber string) error {
	_, existingSO, err := r.isAdminOfSO(adminMobileNumber, soMobileNumber)
	if err != nil {
		return err
	}
	return r.db.Unscoped().Delete(&existingSO).Error
}

func (r repository) ReplaceSO(adminMobileNumber string, oldSOMobileNumber string, newSOMobileNumber string) error {
	existingAdmin, existingOldSO, err := r.isAdminOfSO(adminMobileNumber, oldSOMobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Checking if new so is already registered by current admin")
	newSO, err := utils.GetSOBy(r.db, newSOMobileNumber)
	if err != nil {
		logrus.Info("Adding new so")
		newSO = models.SupervisingOfficer{MobileNumber: newSOMobileNumber, AdminID: existingAdmin.ID}
		err = r.db.Save(&newSO).Error
		if err != nil {
			logrus.Error("Could not save new SO ", err)
			return err
		}
	}
	if newSO.AdminID != existingAdmin.ID {
		logrus.Error(constants.SONotRegisteredByAdmin)
		return errors.New(constants.SONotRegisteredByAdmin)
	}
	quarantine := &models.Quarantine{SupervisingOfficerID: existingOldSO.ID}
	return r.db.Model(&models.Quarantine{}).Where(quarantine).Update(models.Quarantine{SupervisingOfficerID: newSO.ID}).Error
}

func (r repository) isAdminOfSO(adminMobileNumber string, soMobileNumber string) (*models.Admin, *models.SupervisingOfficer, error) {
	existingAdmin, err := r.getBy(adminMobileNumber)
	if err != nil {
		logrus.Error("Admin not found")
		return nil, nil, err
	}
	existingSO, err := utils.GetSOBy(r.db, soMobileNumber)
	if err != nil {
		logrus.Error("SO not found")
		return &existingAdmin, nil, err
	}
	logrus.Info("Checking if so is registered by current admin")
	if existingSO.AdminID != existingAdmin.ID {
		logrus.Error(constants.SONotRegisteredByAdmin)
		return &existingAdmin, &existingSO, errors.New(constants.SONotRegisteredByAdmin)
	}
	return &existingAdmin, &existingSO, nil
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
