package so

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
	AddQuarantine(mobileNumber string, quarantine models.Quarantine) error
	GetQuarantines(mobileNumber string) ([]models.Quarantine, error)
	DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) error
}

type repository struct {
	db *gorm.DB
}

func (r repository) IsExists(mobileNumber string) bool {
	user, err := utils.GetSOBy(r.db, mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) AddQuarantine(mobileNumber string, quarantine models.Quarantine) error {
	existingSO, err := utils.GetSOBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	quarantine.SupervisingOfficerID = existingSO.ID
	return r.db.Save(&quarantine).Error
}

func (r repository) GetQuarantines(mobileNumber string) ([]models.Quarantine, error) {
	return utils.GetQuarantines(r.db, mobileNumber)
}

func (r repository) DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) error {
	_, existingQuarantine, err := r.isSOOfQuarantine(soMobileNumber, quarantineMobileNumber)
	if err != nil {
		return err
	}
	return r.db.Unscoped().Delete(existingQuarantine).Error
}

func (r repository) isSOOfQuarantine(soMobileNumber string, quarantineMobileNumber string) (*models.SupervisingOfficer, *models.Quarantine, error) {
	existingSO, err := utils.GetSOBy(r.db, soMobileNumber)
	if err != nil {
		logrus.Error(constants.SONotExistsError, err)
		return nil, nil, errors.New(constants.SONotExistsError)
	}
	existingQuarantine, err := utils.GetQuarantineBy(r.db, quarantineMobileNumber)
	if err != nil {
		logrus.Error(constants.QuarantineNotExistsError, err)
		return &existingSO, nil, errors.New(constants.QuarantineNotExistsError)
	}
	logrus.Info("Checking if quarantine is registered by current so")
	if existingSO.ID != existingQuarantine.SupervisingOfficerID {
		logrus.Error(constants.QuarantineNotRegisteredBySOError)
		return &existingSO, &existingQuarantine, errors.New(constants.QuarantineNotRegisteredBySOError)
	}
	return &existingSO, &existingQuarantine, nil
}


func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
