package so

import (
	"trackcoro/constants"
	models2 "trackcoro/database/models"
	"trackcoro/models"
	"trackcoro/notify"
	"trackcoro/utils"
)

type Service interface {
	Verify(mobileNumber string) bool
	AddQuarantine(soMobileNumber string, quarantineMobileNumber string) *models.Error
	GetQuarantines(soMobileNumber string) ([]models.QuarantineDetails, *models.Error)
	GetQuarantine(soMobileNumber string, quarantineMobileNumber string) (*models.QuarantineDetails, *models.Error)
	DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) *models.Error
	UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models.Error
	NotifyQuarantines(request models.NotificationRequest, soMobileNumber string) *models.Error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) AddQuarantine(soMobileNumber string, quarantineMobileNumber string) *models.Error {
	return s.repository.AddQuarantine(soMobileNumber, models2.Quarantine{MobileNumber: quarantineMobileNumber})
}

func (s service) GetQuarantines(soMobileNumber string) ([]models.QuarantineDetails, *models.Error) {
	quarantinesFromDB, err := s.repository.GetQuarantines(soMobileNumber)
	if err != nil {
		return nil, err
	}
	return utils.GetMappedQuarantines(quarantinesFromDB), nil
}

func (s service) GetQuarantine(soMobileNumber string, quarantineMobileNumber string) (*models.QuarantineDetails, *models.Error) {
	quarantine, err := s.repository.GetQuarantine(soMobileNumber, quarantineMobileNumber)
	if err != nil {
		return nil, err
	}
	mappedQuarantine := utils.GetMappedQuarantine(*quarantine)
	return &mappedQuarantine, nil
}

func (s service) DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) *models.Error {
	return s.repository.DeleteQuarantine(soMobileNumber, quarantineMobileNumber)
}

func (s service) UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models.Error {
	return s.repository.UpdateDeviceTokenId(mobileNumber, deviceTokenId)
}

func (s service) NotifyQuarantines(request models.NotificationRequest, soMobileNumber string) *models.Error {
	quarantines, err := s.repository.GetQuarantines(soMobileNumber)
	if err != nil {
		return err
	}
	deviceTokenIds := getDeviceTokenIds(quarantines)
	failedTokens, err := notify.SendNotification(deviceTokenIds, map[string]string{
		"type":             request.Type,
		"message":          request.Message,
		"so_mobile_number": soMobileNumber,
	})
	if err != nil {
		return err
	}
	if len(failedTokens) > 0 {
		return &constants.SendNotificationFailedError
	}
	return nil
}

func getDeviceTokenIds(quarantines []models2.Quarantine) []string {
	var deviceIds []string
	for _, quarantine := range quarantines {
		deviceIds = append(deviceIds, quarantine.DeviceTokenId)
	}
	return deviceIds
}

func NewService(repository Repository) Service {
	return service{repository}
}
