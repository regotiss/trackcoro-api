package quarantine

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"time"
	"trackcoro/constants"
	dbmodels "trackcoro/database/models"
	models2 "trackcoro/models"
	"trackcoro/objectstorage"
	"trackcoro/quarantine/models"
	"trackcoro/utils"
)

type Service interface {
	Verify(mobileNumber string) bool
	SaveDetails(request models2.QuarantineDetails) error
	GetDaysStatus(mobileNumber string) (models.DaysStatusResponse, error)
	GetDetails(mobileNumber string) (models2.QuarantineDetails, error)
	UploadPhoto(mobileNumber string, photo multipart.File, photoSize int64) error
}

type service struct {
	repository Repository
}

func (s service) UploadPhoto(mobileNumber string, photo multipart.File, photoSize int64) error {
	photoName := fmt.Sprintf("%s.jpg", mobileNumber)
	logrus.Info("file name", photoName)
	photoContent := make([]byte, photoSize)
	_, err := photo.Read(photoContent)
	if err != nil {
		logrus.Error("Unable to read photo content", err)
		return err
	}
	_, err = objectstorage.PutObject(photoName, photoContent)
	return err
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) SaveDetails(detailsRequest models2.QuarantineDetails) error {
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

func (s service) GetDetails(mobileNumber string) (models2.QuarantineDetails, error) {
	quarantine, err := s.repository.GetDetails(mobileNumber)
	if err != nil {
		return models2.QuarantineDetails{}, err
	}
	return mapFromDBQuarantine(quarantine), nil
}

func mapToDBQuarantine(detailRequest models2.QuarantineDetails) (dbmodels.Quarantine, error) {
	DOB, err := time.Parse(constants.DetailsTimeFormat, detailRequest.DOB)
	if err != nil {
		logrus.Error("Could not parse dob ", err)
		return dbmodels.Quarantine{}, errors.New(constants.TimeParseError)
	}
	QuarantineStartedFrom, err := time.Parse(constants.DetailsTimeFormat, detailRequest.QuarantineStartedFrom)

	if err != nil {
		logrus.Error("Could not parse quarantine started from ", err)
		return dbmodels.Quarantine{}, errors.New(constants.TimeParseError)
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
		DeviceTokenId:          detailRequest.DeviceTokenId,
	}, nil
}

func mapFromDBQuarantine(quarantine dbmodels.Quarantine) models2.QuarantineDetails {
	details := utils.GetMappedQuarantine(quarantine)
	return details
}

func mapToDBTravelHistory(travelHistoryRequest []models.TravelHistory) ([]dbmodels.QuarantineTravelHistory, error) {
	var travelHistory []dbmodels.QuarantineTravelHistory
	for _, history := range travelHistoryRequest {
		visitedDate, err := time.Parse(constants.DetailsTimeFormat, history.VisitDate)
		if err != nil {
			logrus.Error("Could not parse visited date of travel ", history.PlaceVisited, " error-", err)
			return nil, errors.New(constants.TimeParseError)
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


func mapToDBAddress(address models2.Address) dbmodels.QuarantineAddress {
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


func NewService(repository Repository) Service {
	return service{repository}
}
