package admin

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/constants"
	dbmodels "trackcoro/database/models"
	"trackcoro/models"
)

type Service interface {
	Verify(mobileNumber string) bool
	Add() error
	AddSO(adminMobileNumber string, soRequest models.SODetails) error
	GetSOs(adminMobileNumber string) ([]models.SODetails, error)
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

func (s service) AddSO(adminMobileNumber string, soRequest models.SODetails) error {
	return s.repository.AddSO(adminMobileNumber, mapToDbSO(soRequest))
}

func (s service) GetSOs(adminMobileNumber string) ([]models.SODetails, error) {
	SOsFromDB, err := s.repository.GetSOs(adminMobileNumber)
	if err != nil {
		return nil, err
	}
	var SOs []models.SODetails
	for _, SO := range SOsFromDB {
		SOs = append(SOs, mapFromDbSO(SO))
	}
	return SOs, nil
}

func mapToDbSO(soRequest models.SODetails) dbmodels.SupervisingOfficer {
	return dbmodels.SupervisingOfficer{
		MobileNumber:         soRequest.MobileNumber,
		Name:                 soRequest.Name,
		BadgeId:              soRequest.BadgeId,
		Designation:          soRequest.Designation,
		PoliceStationAddress: soRequest.PoliceStationAddress,
	}
}

func mapFromDbSO(SO dbmodels.SupervisingOfficer) models.SODetails {
	return models.SODetails{
		MobileNumber:         SO.MobileNumber,
		Name:                 SO.Name,
		BadgeId:              SO.BadgeId,
		Designation:          SO.Designation,
		PoliceStationAddress: SO.PoliceStationAddress,
	}
}

func NewService(repository Repository) Service {
	return service{repository}
}
