package so

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
	"trackcoro/constants"
	"trackcoro/database/models"
	models2 "trackcoro/models"
	"trackcoro/utils"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	AddQuarantine(mobileNumber string, quarantine models.Quarantine) *models2.Error
	GetQuarantines(mobileNumber string) ([]models.Quarantine, *models2.Error)
	GetQuarantine(soMobileNumber string, quarantineMobileNumber string) (*models.Quarantine, *models2.Error)
	DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) *models2.Error
	UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models2.Error
	SaveUploadDetails(quarantineMobileNumber string) *models2.Error
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

func (r repository) AddQuarantine(mobileNumber string, quarantine models.Quarantine) *models2.Error {
	existingSO, err := utils.GetSOBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	quarantine.SupervisingOfficerID = existingSO.ID
	dbError := r.db.Save(&quarantine).Error
	if dbError != nil {
		logrus.Error("Couldn't save quarantine ", dbError)
		return &constants.QuarantineAlreadyExistsError
	}
	return nil
}

func (r repository) GetQuarantines(mobileNumber string) ([]models.Quarantine, *models2.Error) {
	return utils.GetQuarantinesForSO(r.db, mobileNumber)
}

func (r repository) GetQuarantine(soMobileNumber string, quarantineMobileNumber string) (*models.Quarantine, *models2.Error) {
	_, quarantine, err := r.isSOOfQuarantine(soMobileNumber, quarantineMobileNumber)
	if err != nil {
		return nil, err
	}
	r.db.Preload("Address").Preload("TravelHistory").Preload("PhotoUpload").Where(&quarantine).First(&quarantine)
	return quarantine, nil
}

func (r repository) DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) *models2.Error {
	_, existingQuarantine, err := r.isSOOfQuarantine(soMobileNumber, quarantineMobileNumber)
	if err != nil {
		return err
	}
	dbError := r.db.Delete(existingQuarantine).Error
	if dbError != nil {
		logrus.Error("Couldn't delete quarantine ", dbError)
		return &constants.InternalError
	}
	return nil
}

func (r repository) UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models2.Error {
	user, err := utils.GetSOBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Updating device token id")
	userWithTokenId := &models.SupervisingOfficer{DeviceTokenId: deviceTokenId}
	dbError := r.db.Model(&user).Update(userWithTokenId).Error
	if dbError != nil {
		logrus.Error("Could not save device token id ", dbError)
		return &constants.InternalError
	}
	return nil
}

func (r repository) SaveUploadDetails(mobileNumber string) *models2.Error {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return &constants.QuarantineNotExistsError
	}
	photoUpload := models.PhotoUpload{QuarantineID: user.ID}
	dbError := r.db.Where(photoUpload).First(&photoUpload).Error
	if dbError != nil {
		photoUpload.RequestedOn = time.Now()
		dbError := r.db.Create(&photoUpload).Error
		if dbError != nil {
			return &constants.InternalError
		}
	}
	photoUpload.RequestedOn = time.Now()
	dbError = r.db.Update(&photoUpload).Error
	if dbError != nil {
		return &constants.InternalError
	}
	return nil
}

func (r repository) isSOOfQuarantine(soMobileNumber string, quarantineMobileNumber string) (*models.SupervisingOfficer, *models.Quarantine, *models2.Error) {
	existingSO, err := utils.GetSOBy(r.db, soMobileNumber)
	if err != nil {
		return nil, nil, &constants.SONotExistsError
	}
	existingQuarantine, quaError := utils.GetQuarantineBy(r.db, quarantineMobileNumber)
	if quaError != nil {
		return &existingSO, nil, &constants.QuarantineNotExistsError
	}
	logrus.Info("Checking if quarantine is registered by current so")
	if existingSO.ID != existingQuarantine.SupervisingOfficerID {
		logrus.Error("quarantine is not registered by so")
		return &existingSO, &existingQuarantine, &constants.QuarantineNotRegisteredBySOError
	}
	return &existingSO, &existingQuarantine, nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
