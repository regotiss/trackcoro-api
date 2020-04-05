package quarantine

import (
	"errors"
	"github.com/sirupsen/logrus"
	"time"
	"trackcoro/quarantine/models"
)

type Service interface {
	Verify(mobileNumber string) bool
	SaveDetails(request models.SaveDetailsRequest) error
	GetDaysStatus(mobileNumber string) (models.DaysStatusResponse, error)
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.isExists(mobileNumber)
}

func (s service) SaveDetails(detailsRequest models.SaveDetailsRequest) error {
	user, err := mapQuarantine(detailsRequest)
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

func mapQuarantine(detailRequest models.SaveDetailsRequest) (models.Quarantine, error) {
	DOB, err := time.Parse(DetailsTimeFormat, detailRequest.DOB)
	if err != nil {
		logrus.Error("Could not parse dob ", err)
		return models.Quarantine{}, errors.New(TimeParseError)
	}
	QuarantineStartedFrom, err := time.Parse(DetailsTimeFormat, detailRequest.QuarantineStartedFrom)

	if err != nil {
		logrus.Error("Could not parse quarantine started from ", err)
		return models.Quarantine{}, errors.New(TimeParseError)
	}
	history, err := mapTravelHistory(detailRequest.TravelHistory)
	if err != nil {
		return models.Quarantine{}, err
	}
	return models.Quarantine{
		MobileNumber:           detailRequest.MobileNumber,
		Name:                   detailRequest.Name,
		Address:                mapAddress(detailRequest.Address),
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

func mapTravelHistory(travelHistoryRequest []models.TravelHistory) ([]models.QuarantineTravelHistory, error) {
	var travelHistory []models.QuarantineTravelHistory
	for _, history := range travelHistoryRequest {
		visitedDate, err := time.Parse(DetailsTimeFormat, history.VisitDate)
		if err != nil {
			logrus.Error("Could not parse visited date of travel ", history.PlaceVisited, " error-", err)
			return nil, errors.New(TimeParseError)
		}
		travelHistory = append(travelHistory, models.QuarantineTravelHistory{
			PlaceVisited:         history.PlaceVisited,
			VisitDate:            visitedDate,
			TimeSpentInDays:      history.TimeSpentInDays,
			ModeOfTransportation: history.ModeOfTransportation,
		})
	}
	return travelHistory, nil
}

func mapAddress(address models.Address) models.QuarantineAddress {
	return models.QuarantineAddress{
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

func NewService(repository Repository) Service {
	return service{repository}
}
