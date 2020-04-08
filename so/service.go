package so

import (
	"trackcoro/constants"
	models2 "trackcoro/database/models"
	"trackcoro/models"
	"trackcoro/utils"
)

type Service interface {
	Verify(mobileNumber string) bool
	AddQuarantine(soMobileNumber string, quarantineMobileNumber string) error
	GetQuarantines(soMobileNumber string) ([]models.QuarantineDetails, error)
	DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) AddQuarantine(soMobileNumber string, quarantineMobileNumber string) error {
	return s.repository.AddQuarantine(soMobileNumber, models2.Quarantine{MobileNumber: quarantineMobileNumber})
}

func (s service) GetQuarantines(soMobileNumber string) ([]models.QuarantineDetails, error) {
	quarantinesFromDB, err := s.repository.GetQuarantines(soMobileNumber)
	if err != nil {
		return nil, err
	}
	var quarantineDetails []models.QuarantineDetails
	for _, quarantine := range quarantinesFromDB {
		quarantineDetails = append(quarantineDetails, getQuarantineDetails(quarantine))
	}
	return quarantineDetails, nil
}

func getQuarantineDetails(quarantine models2.Quarantine) models.QuarantineDetails {
	mappedQuarantine := utils.GetMappedQuarantine(quarantine)
	if quarantine.CurrentLocationLatitude != constants.Empty || quarantine.CurrentLocationLongitude != constants.Empty {
		mappedQuarantine.CurrentLocation = &models.Coordinates{
			Latitude:  quarantine.CurrentLocationLatitude,
			Longitude: quarantine.CurrentLocationLongitude,
		}
	}
	return mappedQuarantine
}

func (s service) DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) error {
	return s.repository.DeleteQuarantine(soMobileNumber, quarantineMobileNumber)
}


func NewService(repository Repository) Service {
	return service{repository}
}
