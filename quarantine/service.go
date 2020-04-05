package quarantine

import (
	"errors"
	"github.com/sirupsen/logrus"
	"time"
	dbmodels "trackcoro/database/models"
	"trackcoro/quarantine/models"
)

type Service interface {
	Verify(mobileNumber string) bool
	SaveDetails(request models.ProfileDetails) error
	GetDaysStatus(mobileNumber string) (models.DaysStatusResponse, error)
	GetDetails(mobileNumber string) (models.ProfileDetails, error)
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.isExists(mobileNumber)
}

func (s service) SaveDetails(detailsRequest models.ProfileDetails) error {
	user, err := mapToDBQuarantine(detailsRequest)
	if err != nil {
		return err
	}
	return s.repository.SaveDetails(user)
}

func (s service) GetDaysStatus(mobileNumber string) (models.DaysStatusResponse, error) {
	days, startedFrom, err := s.repository.GetQuarantineDays(mobileNumber)
	if err != nil {
		return models.DaysStatusResponse{}, err
	}
	remainingDays := int(days - uint(time.Now().Sub(startedFrom).Hours()/24))
	if remainingDays < 0 {
		remainingDays = 0
	}
	daysStatus := models.DaysStatusResponse{
		NoOfQuarantineDays: days,
		StatedFrom:         startedFrom,
		RemainingDays:      remainingDays,
	}
	return daysStatus, nil
}

func (s service) GetDetails(mobileNumber string) (models.ProfileDetails, error) {
	quarantine, err := s.repository.GetDetails(mobileNumber)
	if err != nil {
		return models.ProfileDetails{}, err
	}
	return mapFromDBQuarantine(quarantine), nil
}

func mapToDBQuarantine(detailRequest models.ProfileDetails) (dbmodels.Quarantine, error) {
	DOB, err := time.Parse(DetailsTimeFormat, detailRequest.DOB)
	if err != nil {
		logrus.Error("Could not parse dob ", err)
		return dbmodels.Quarantine{}, errors.New(TimeParseError)
	}
	QuarantineStartedFrom, err := time.Parse(DetailsTimeFormat, detailRequest.QuarantineStartedFrom)

	if err != nil {
		logrus.Error("Could not parse quarantine started from ", err)
		return dbmodels.Quarantine{}, errors.New(TimeParseError)
	}
	history, err := mapToDBTravelHistory(detailRequest.TravelHistory)
	if err != nil {
		return dbmodels.Quarantine{}, err
	}
	return dbmodels.Quarantine{
		MobileNumber:           detailRequest.MobileNumber,
		Name:                   detailRequest.Name,
		Address:                mapToDBAddress(detailRequest.Address),
		Occupation:             detailRequest.Occupation,
		DOB:                    DOB,
		AnyPractitionerConsult: detailRequest.AnyPractitionerConsult,
		NoOfQuarantineDays:     detailRequest.NoOfQuarantineDays,
		QuarantineStartedFrom:  QuarantineStartedFrom,
		FamilyMembers:          detailRequest.FamilyMembers,
		SecondaryContactNumber: detailRequest.SecondaryContactNumber,
		TravelHistory:          history,
	}, nil
}

func mapFromDBQuarantine(quarantine dbmodels.Quarantine) models.ProfileDetails {
	return models.ProfileDetails{
		MobileNumber:           quarantine.MobileNumber,
		Name:                   quarantine.Name,
		Address:                mapFromDBAddress(quarantine.Address),
		Occupation:             quarantine.Occupation,
		DOB:                    quarantine.DOB.String(),
		TravelHistory:          mapFromDBTravelHistory(quarantine.TravelHistory),
		AnyPractitionerConsult: quarantine.AnyPractitionerConsult,
		NoOfQuarantineDays:     quarantine.NoOfQuarantineDays,
		QuarantineStartedFrom:  quarantine.QuarantineStartedFrom.String(),
		FamilyMembers:          quarantine.FamilyMembers,
		SecondaryContactNumber: quarantine.SecondaryContactNumber,
	}
}
func mapToDBTravelHistory(travelHistoryRequest []models.TravelHistory) ([]dbmodels.QuarantineTravelHistory, error) {
	var travelHistory []dbmodels.QuarantineTravelHistory
	for _, history := range travelHistoryRequest {
		visitedDate, err := time.Parse(DetailsTimeFormat, history.VisitDate)
		if err != nil {
			logrus.Error("Could not parse visited date of travel ", history.PlaceVisited, " error-", err)
			return nil, errors.New(TimeParseError)
		}
		travelHistory = append(travelHistory, dbmodels.QuarantineTravelHistory{
			PlaceVisited:         history.PlaceVisited,
			VisitDate:            visitedDate,
			TimeSpentInDays:      history.TimeSpentInDays,
			ModeOfTransportation: history.ModeOfTransportation,
		})
	}
	return travelHistory, nil
}

func mapFromDBTravelHistory(quarantineTravelHistory []dbmodels.QuarantineTravelHistory) []models.TravelHistory{
	var travelHistory []models.TravelHistory
	for _, history := range quarantineTravelHistory {
		travelHistory = append(travelHistory, models.TravelHistory{
			PlaceVisited:         history.PlaceVisited,
			VisitDate:            history.VisitDate.String(),
			TimeSpentInDays:      history.TimeSpentInDays,
			ModeOfTransportation: history.ModeOfTransportation,
		})
	}
	return travelHistory
}

func mapToDBAddress(address models.Address) dbmodels.QuarantineAddress {
	return dbmodels.QuarantineAddress{
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		AddressLine3: address.AddressLine3,
		Locality:     address.Locality,
		City:         address.City,
		District:     address.District,
		State:        address.State,
		Country:      address.Country,
		PinCode:      address.PinCode,
		Latitude:     address.Coordinates.Latitude,
		Longitude:    address.Coordinates.Longitude,
	}
}

func mapFromDBAddress(address dbmodels.QuarantineAddress) models.Address {
	return models.Address {
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		AddressLine3: address.AddressLine3,
		Locality:     address.Locality,
		City:         address.City,
		District:     address.District,
		State:        address.State,
		Country:      address.Country,
		PinCode:      address.PinCode,
		Coordinates:  models.Coordinates{
			Latitude:  address.Latitude,
			Longitude: address.Longitude,
		},
	}
}

func NewService(repository Repository) Service {
	return service{repository}
}
