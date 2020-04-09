package admin

import (
	"github.com/sirupsen/logrus"
	"os"
	"trackcoro/constants"
	dbmodels "trackcoro/database/models"
	"trackcoro/models"
	"trackcoro/utils"
)

type Service interface {
	Verify(mobileNumber string) bool
	Add() *models.Error
	AddSO(adminMobileNumber string, soRequest models.SODetails) *models.Error
	GetSOs(adminMobileNumber string) ([]models.SODetails, *models.Error)
	GetQuarantines(adminMobileNumber string, soMobileNumber string) ([]models.QuarantineDetails, *models.Error)
	DeleteSO(adminMobileNumber string, soMobileNumber string) *models.Error
	ReplaceSO(adminMobileNumber string, oldSOMobileNumber string, newSOMobileNumber string) *models.Error
	DeleteAllSOs(adminMobileNumber string) *models.Error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) Add() *models.Error {
	mobileNumber := os.Getenv(constants.AdminMobileNumber)
	if mobileNumber == constants.Empty {
		logrus.Error(constants.EnvVariableNotFoundError)
		return &constants.InternalError
	}
	return s.repository.Add(dbmodels.Admin{MobileNumber:mobileNumber})
}

func (s service) AddSO(adminMobileNumber string, soRequest models.SODetails) *models.Error {
	return s.repository.AddSO(adminMobileNumber, mapToDbSO(soRequest))
}

func (s service) GetSOs(adminMobileNumber string) ([]models.SODetails, *models.Error) {
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

func (s service) GetQuarantines(adminMobileNumber string, soMobileNumber string) ([]models.QuarantineDetails, *models.Error) {
	quarantinesFromDB, err := s.repository.GetQuarantines(adminMobileNumber, soMobileNumber)
	if err != nil {
		return nil, err
	}

	return utils.GetMappedQuarantines(quarantinesFromDB), nil
}

func (s service) DeleteSO(adminMobileNumber string, soMobileNumber string) *models.Error {
	return s.repository.DeleteSO(adminMobileNumber, soMobileNumber)
}

func (s service) ReplaceSO(adminMobileNumber string, oldSOMobileNumber string, newSOMobileNumber string) *models.Error {
	return s.repository.ReplaceSO(adminMobileNumber, oldSOMobileNumber, newSOMobileNumber)
}

func (s service) DeleteAllSOs(adminMobileNumber string) *models.Error {
	return s.repository.DeleteAllSOs(adminMobileNumber)
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
