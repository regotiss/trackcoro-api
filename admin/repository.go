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
	DeleteAllSOs(adminMobileNumber string) error
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
	err = r.db.Save(&admin).Error
	if err != nil {
		return  err
	}
	return r.db.Save(&models.SupervisingOfficer{MobileNumber: admin.MobileNumber, AdminID: admin.ID}).Error
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
	transaction := r.db.Begin()

	logrus.Info("Assigning admin as a so to quarantines of deleting so")
	adminSO, err := utils.GetSOBy(r.db, adminMobileNumber)
	if err != nil {
		logrus.Error("Could not assign admin so")
		return err
	}
	quarantine := &models.Quarantine{SupervisingOfficerID: existingSO.ID}
	err = r.db.Model(&models.Quarantine{}).Where(quarantine).Update(models.Quarantine{SupervisingOfficerID: adminSO.ID}).Error
	if err != nil {
		transaction.Rollback()
		return err
	}
	err = transaction.Unscoped().Delete(&existingSO).Error
	if err != nil {
		transaction.Rollback()
		return err
	}

	return transaction.Commit().Error
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
		logrus.Error(constants.SONotRegisteredByAdminError)
		return errors.New(constants.SONotRegisteredByAdminError)
	}
	quarantine := &models.Quarantine{SupervisingOfficerID: existingOldSO.ID}
	return r.db.Model(&models.Quarantine{}).Where(quarantine).Update(models.Quarantine{SupervisingOfficerID: newSO.ID}).Error
}

func (r repository) DeleteAllSOs(adminMobileNumber string) error {
	r.db.LogMode(true)
	existingAdmin, err := r.getBy(adminMobileNumber)
	if err != nil {
		return err
	}

	sos := &models.SupervisingOfficer{AdminID: existingAdmin.ID}
	return r.db.Unscoped().Where(sos).Delete(sos).Error
}

func (r repository) isAdminOfSO(adminMobileNumber string, soMobileNumber string) (*models.Admin, *models.SupervisingOfficer, error) {
	existingAdmin, err := r.getBy(adminMobileNumber)
	if err != nil {
		logrus.Error(constants.AdminNotExistsError, err)
		return nil, nil, err
	}
	existingSO, err := utils.GetSOBy(r.db, soMobileNumber)
	if err != nil {
		logrus.Error(err)
		return &existingAdmin, nil, errors.New(constants.SONotExistsError)
	}
	logrus.Info("Checking if so is registered by current admin")
	if existingSO.AdminID != existingAdmin.ID {
		logrus.Error(constants.SONotRegisteredByAdminError)
		return &existingAdmin, &existingSO, errors.New(constants.SONotRegisteredByAdminError)
	}
	return &existingAdmin, &existingSO, nil
}

func (r repository) getBy(mobileNumber string) (models.Admin, error) {
	var user models.Admin
	err := r.db.Where(&models.Admin{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return models.Admin{}, errors.New(constants.AdminNotExistsError)
	}
	return user, nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
