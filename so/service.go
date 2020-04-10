package so

import (
	models2 "trackcoro/database/models"
	"trackcoro/models"
	"trackcoro/utils"
)

type Service interface {
	Verify(mobileNumber string) bool
	AddQuarantine(soMobileNumber string, quarantineMobileNumber string) error
	GetQuarantines(soMobileNumber string) ([]models.QuarantineDetails, *models.Error)
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

func (s service) GetQuarantines(soMobileNumber string) ([]models.QuarantineDetails, *models.Error) {
	quarantinesFromDB, err := s.repository.GetQuarantines(soMobileNumber)
	if err != nil {
		return nil, err
	}
	return utils.GetMappedQuarantines(quarantinesFromDB), nil
}

func (s service) DeleteQuarantine(soMobileNumber string, quarantineMobileNumber string) error {
	return s.repository.DeleteQuarantine(soMobileNumber, quarantineMobileNumber)
}

func NewService(repository Repository) Service {
	return service{repository}
}
