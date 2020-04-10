package so

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"trackcoro/constants"
	"trackcoro/database/models"
	models2 "trackcoro/models"
	"trackcoro/utils"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	AddQuarantine(mobileNumber string, quarantine models.Quarantine) error
	GetQuarantines(mobileNumber string) ([]models.Quarantine, *models2.Error)
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

func (r repository) GetQuarantines(mobileNumber string) ([]models.Quarantine, *models2.Error) {
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
		return nil, nil, constants.SONotExistsError
	}
	existingQuarantine, quaError := utils.GetQuarantineBy(r.db, quarantineMobileNumber)
	if quaError != nil {
		return &existingSO, nil, constants.QuarantineNotExistsError
	}
	logrus.Info("Checking if quarantine is registered by current so")
	if existingSO.ID != existingQuarantine.SupervisingOfficerID {
		logrus.Error("quarantine is not registered by so")
		return &existingSO, &existingQuarantine, errors.New(constants.QuarantineNotRegisteredBySOError)
	}
	return &existingSO, &existingQuarantine, nil
}


func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
