package quarantine

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
	IsExists(mobileNumber string) (bool, bool)
	SaveDetails(quarantine models.Quarantine) *models2.Error
	GetQuarantineDays(mobileNumber string) (uint, time.Time, *models2.Error)
	GetDetails(mobileNumber string) (models.Quarantine, *models2.Error)
	UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng string) *models2.Error
	UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models2.Error
	SaveUploadDetails(mobileNumber string) *models2.Error
}

type repository struct {
	db *gorm.DB
}

func (r repository) IsExists(mobileNumber string) (bool, bool) {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return false, false
	}
	return user.MobileNumber == mobileNumber, user.Name != constants.Empty
}

func (r repository) SaveDetails(quarantine models.Quarantine) *models2.Error {
	user, err := utils.GetQuarantineBy(r.db, quarantine.MobileNumber)
	if err != nil {
		return err
	}
	quarantine.ID = user.ID
	quarantine.CurrentLocationLatitude = quarantine.Address.Latitude
	quarantine.CurrentLocationLongitude = quarantine.Address.Longitude
	quarantine.SupervisingOfficerID = user.SupervisingOfficerID
	r.db.Unscoped().Delete(models.QuarantineAddress{QuarantineID: user.ID})
	r.db.Unscoped().Delete(models.QuarantineTravelHistory{QuarantineID: user.ID})
	dbErr := r.db.Save(&quarantine).Error
	if dbErr != nil {
		logrus.Error("Could not save details ", dbErr)
		return &constants.InternalError
	}
	return nil
}

func (r repository) GetQuarantineDays(mobileNumber string) (uint, time.Time, *models2.Error) {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return 0, time.Time{}, err
	}
	return user.NoOfQuarantineDays, user.QuarantineStartedFrom, nil
}

func (r repository) GetDetails(mobileNumber string) (models.Quarantine, *models2.Error) {
	user, err := utils.GetAllQuarantineDetails(r.db, mobileNumber)
	if err != nil {
		return models.Quarantine{}, err
	}
	return user, nil
}

func (r repository) UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng string) *models2.Error {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Updating current location")
	userWithCurrentLocation := &models.Quarantine{
		CurrentLocationLatitude:  currentLocationLat,
		CurrentLocationLongitude: currentLocationLng,
	}
	dbError := r.db.Model(&user).Update(userWithCurrentLocation).Error
	if dbError != nil {
		logrus.Error("Could not save current location ", dbError)
		return &constants.InternalError
	}
	return nil
}

func (r repository) UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models2.Error {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Updating device token id")
	userWithTokenId := &models.Quarantine{DeviceTokenId: deviceTokenId}
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
		photoUpload.UploadedOn = time.Now()
		dbError := r.db.Create(photoUpload).Error
		if dbError != nil {
			return &constants.InternalError
		}
	}
	photoUpload.UploadedOn = time.Now()
	dbError = r.db.Update(&photoUpload).Error
	if dbError != nil {
		return &constants.InternalError
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
