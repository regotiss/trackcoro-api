package admin

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
	models2 "trackcoro/models"
	"trackcoro/utils"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	Add(admin models.Admin) *models2.Error
	AddSO(adminMobileNumber string, so models.SupervisingOfficer) *models2.Error
	GetSOs(adminMobileNumber string) ([]models.SupervisingOfficer, *models2.Error)
	GetQuarantines(adminMobileNumber string, soMobileNumber string) ([]models.Quarantine, *models2.Error)
	DeleteSO(adminMobileNumber string, soMobileNumber string) *models2.Error
	ReplaceSO(adminMobileNumber string, oldSOMobileNumber string, newSOMobileNumber string) *models2.Error
	DeleteAllSOs(adminMobileNumber string) *models2.Error
}

type repository struct {
	db *gorm.DB
}

func (r repository) IsExists(mobileNumber string) bool {
	user, err := r.getAdminBy(mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) Add(admin models.Admin) *models2.Error {
	_, err := r.getAdminBy(admin.MobileNumber)
	if err == nil {
		logrus.Info("Admin already exists")
		return nil
	}

	dbErr := r.db.Create(&admin).Error
	if dbErr != nil {
		return &constants.InternalError
	}
	if dbErr = r.db.Save(&models.SupervisingOfficer{MobileNumber: admin.MobileNumber, AdminID: admin.ID}).Error; dbErr != nil {
		return &constants.InternalError
	}
	return nil
}

func (r repository) AddSO(mobileNumber string, so models.SupervisingOfficer) *models2.Error {
	existingAdmin, err := r.getAdminBy(mobileNumber)
	if err != nil {
		return err
	}
	so.AdminID = existingAdmin.ID
	dbError := r.db.Save(&so).Error
	if dbError != nil {
		logrus.Error("Could  not save SO ", dbError.Error())
		return &constants.SOAlreadyExistsError
	}
	return nil
}

func (r repository) GetSOs(adminMobileNumber string) ([]models.SupervisingOfficer, *models2.Error) {
	existingAdmin, err := r.getAdminBy(adminMobileNumber)
	if err != nil {
		return nil, err
	}
	var SOs []models.SupervisingOfficer
	dbError := r.db.Model(&existingAdmin).Where("mobile_number <> ?", adminMobileNumber).Related(&SOs).Error
	if dbError != nil {
		logrus.Error("Could not fetch SOs ", dbError.Error())
		return SOs, &constants.InternalError
	}
	return SOs, nil
}

func (r repository) GetQuarantines(adminMobileNumber string, soMobileNumber string) ([]models.Quarantine, *models2.Error) {
	_, _, err := r.isAdminOfSO(adminMobileNumber, soMobileNumber)
	if err != nil {
		return nil, err
	}
	return utils.GetQuarantines(r.db, soMobileNumber)
}

func (r repository) DeleteSO(adminMobileNumber string, soMobileNumber string) *models2.Error {
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
	dbError := r.db.Model(&models.Quarantine{}).Where(quarantine).Update(models.Quarantine{SupervisingOfficerID: adminSO.ID}).Error
	if dbError != nil {
		logrus.Error("Could not assign quarantine so as admin so ", dbError.Error())
		transaction.Rollback()
		return &constants.InternalError
	}
	dbError = transaction.Unscoped().Delete(&existingSO).Error
	if dbError != nil {
		logrus.Error("Could not delete so ", dbError.Error())
		transaction.Rollback()
		return &constants.InternalError
	}

	dbError = transaction.Commit().Error
	if dbError != nil {
		logrus.Error("Could not commit changes ", dbError.Error())
		return &constants.InternalError
	}
	return nil
}

func (r repository) ReplaceSO(adminMobileNumber string, oldSOMobileNumber string, newSOMobileNumber string) *models2.Error {
	existingAdmin, existingOldSO, err := r.isAdminOfSO(adminMobileNumber, oldSOMobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Checking if new so is already registered by current admin")
	newSO, err := utils.GetSOBy(r.db, newSOMobileNumber)
	if err != nil {
		return err
	}
	if newSO.AdminID != existingAdmin.ID {
		logrus.Error("Replacing SO not registered by admin")
		return &constants.SONotRegisteredByAdminError
	}

	quarantine := &models.Quarantine{SupervisingOfficerID: existingOldSO.ID}
	dbError := r.db.Model(&models.Quarantine{}).Where(quarantine).Update(models.Quarantine{SupervisingOfficerID: newSO.ID}).Error
	if dbError != nil {
		logrus.Error("Could not replace so ", dbError.Error())
		return &constants.InternalError
	}
	return nil
}

func (r repository) DeleteAllSOs(adminMobileNumber string) *models2.Error {
	r.db.LogMode(true)
	existingAdmin, err := r.getAdminBy(adminMobileNumber)
	if err != nil {
		return err
	}

	sos := &models.SupervisingOfficer{AdminID: existingAdmin.ID}
	dbError := r.db.Unscoped().Where(sos).Delete(sos).Error
	if dbError != nil {
		logrus.Error("Could not delete SOs ", dbError.Error())
		return &constants.InternalError
	}
	return nil
}

func (r repository) isAdminOfSO(adminMobileNumber string, soMobileNumber string) (*models.Admin, *models.SupervisingOfficer, *models2.Error) {
	logrus.Info("Checking if admin exists")
	existingAdmin, err := r.getAdminBy(adminMobileNumber)
	if err != nil {
		return nil, nil, err
	}
	logrus.Info("Checking if so exists")
	existingSO, err := utils.GetSOBy(r.db, soMobileNumber)
	if err != nil {
		return &existingAdmin, nil, err
	}
	logrus.Info("Checking if so is registered by current admin")
	if existingSO.AdminID != existingAdmin.ID {
		logrus.Error(constants.SONotRegisteredByAdminError)
		return &existingAdmin, &existingSO, &constants.SONotRegisteredByAdminError
	}
	return &existingAdmin, &existingSO, nil
}

func (r repository) getAdminBy(mobileNumber string) (models.Admin, *models2.Error) {
	var user models.Admin
	err := r.db.Where(&models.Admin{MobileNumber: mobileNumber}).First(&user).Error
	if err != nil {
		logrus.Error("Could not check mobile number in db ", err)
		return models.Admin{}, &constants.AdminNotExistsError
	}
	return user, nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
