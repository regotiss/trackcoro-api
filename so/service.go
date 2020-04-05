package so

import (
	models2 "trackcoro/database/models"
)

type Service interface {
	Verify(mobileNumber string) bool
	AddQuarantine(soMobileNumber string, quarantineMobileNumber string) error
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

func NewService(repository Repository) Service {
	return service{repository}
}
