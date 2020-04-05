package admin

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/admin/models"
	"trackcoro/constants"
	dbmodels "trackcoro/database/models"
)

type Service interface {
	Verify(mobileNumber string) bool
	Add() error
	AddSO(adminMobileNumber string, soRequest models.AddSORequest) error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) Add() error {
	mobileNumber := os.Getenv(constants.AdminMobileNumber)
	if mobileNumber == constants.Empty {
		logrus.Error(constants.EnvVariableNotFound)
		return errors.New(constants.EnvVariableNotFound)
	}
	return s.repository.Add(dbmodels.Admin{MobileNumber:mobileNumber})
}

func (s service) AddSO(adminMobileNumber string, soRequest models.AddSORequest) error {
	return s.repository.AddSO(adminMobileNumber, mapToSO(soRequest))
}

func mapToSO(soRequest models.AddSORequest) dbmodels.SupervisingOfficer {
	return dbmodels.SupervisingOfficer{
		MobileNumber:         soRequest.MobileNumber,
		Name:                 soRequest.Name,
		BadgeId:              soRequest.BadgeId,
		Designation:          soRequest.Designation,
		PoliceStationAddress: soRequest.PoliceStationAddress,
	}
}
func NewService(repository Repository) Service {
	return service{repository}
}
