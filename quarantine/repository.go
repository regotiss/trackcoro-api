package quarantine

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
	"trackcoro/database/models"
	"trackcoro/utils"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	SaveDetails(quarantine models.Quarantine) error
	GetQuarantineDays(mobileNumber string) (uint, time.Time, error)
	GetDetails(mobileNumber string) (models.Quarantine, error)
	UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng string) error
	UpdateDeviceTokenId(mobileNumber, deviceTokenId string) error
}

type repository struct {
	db *gorm.DB
}

func (r repository) IsExists(mobileNumber string) bool {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) SaveDetails(quarantine models.Quarantine) error {
	user, err := utils.GetQuarantineBy(r.db, quarantine.MobileNumber)
	if err != nil {
		return err
	}
	quarantine.ID = user.ID
	quarantine.SupervisingOfficerID = user.SupervisingOfficerID
	r.db.Unscoped().Delete(models.QuarantineAddress{QuarantineID: user.ID})
	r.db.Unscoped().Delete(models.QuarantineTravelHistory{QuarantineID: user.ID})
	err = r.db.Save(&quarantine).Error
	if err != nil {
		logrus.Error("Could not save details ", err)
	}
	return err
}

func (r repository) GetQuarantineDays(mobileNumber string) (uint, time.Time, error) {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return 0, time.Time{}, err
	}
	return user.NoOfQuarantineDays, user.QuarantineStartedFrom, nil
}

func (r repository) GetDetails(mobileNumber string) (models.Quarantine, error) {
	user, err := utils.GetAllQuarantineDetails(r.db, mobileNumber)
	if err != nil {
		return models.Quarantine{}, err
	}
	return user, nil
}

func (r repository) UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng string) error {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Updating current location")
	userWithCurrentLocation := &models.Quarantine{
		CurrentLocationLatitude: currentLocationLat,
		CurrentLocationLongitude: currentLocationLng,
	}
	return r.db.Model(&user).Update(userWithCurrentLocation).Error
}

func (r repository) UpdateDeviceTokenId(mobileNumber, deviceTokenId string) error {
	user, err := utils.GetQuarantineBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	logrus.Info("Updating device token id")
	userWithTokenId := &models.Quarantine{ DeviceTokenId: deviceTokenId }
	return r.db.Model(&user).Update(userWithTokenId).Error
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
