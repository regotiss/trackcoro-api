package quarantine

import (
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"time"
	"trackcoro/constants"
	dbmodels "trackcoro/database/models"
	models2 "trackcoro/models"
	"trackcoro/notify"
	"trackcoro/objectstorage"
	"trackcoro/quarantine/models"
	"trackcoro/utils"
)

type Service interface {
	Verify(mobileNumber string) bool
	SaveDetails(request models2.QuarantineDetails) *models2.Error
	GetDaysStatus(mobileNumber string) (models.DaysStatusResponse, *models2.Error)
	GetDetails(mobileNumber string) (models2.QuarantineDetails, *models2.Error)
	UploadPhoto(mobileNumber string, photo multipart.File, fileHeader *multipart.FileHeader) *models2.Error
	DownloadPhoto(mobileNumber string) ([]byte, *models2.Error)
	UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng string) *models2.Error
	UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models2.Error
	NotifySO(request models2.NotificationRequest, mobileNumber string) *models2.Error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func (s service) SaveDetails(detailsRequest models2.QuarantineDetails) *models2.Error {
	user, err := mapToDBQuarantine(detailsRequest)
	if err != nil {
		return err
	}
	return s.repository.SaveDetails(user)
}

func (s service) GetDaysStatus(mobileNumber string) (models.DaysStatusResponse, *models2.Error) {
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

func (s service) UploadPhoto(mobileNumber string, photo multipart.File, fileHeader *multipart.FileHeader) *models2.Error {
	photoContent := make([]byte, fileHeader.Size)

	logrus.Info("Reading content of file")
	_, err := photo.Read(photoContent)
	if err != nil {
		logrus.Error("Unable to read photo content", err)
		return &constants.UploadFileContentReadError
	}

	logrus.Info("Uploading file")
	_, err = objectstorage.PutObject(mobileNumber, photoContent)

	if err != nil {
		return &constants.UploadFileFailureError
	}
	return nil
}

func (s service) DownloadPhoto(mobileNumber string) ([]byte, *models2.Error) {
	content, err := objectstorage.GetObject(mobileNumber)
	if err != nil {
		return nil, &constants.DownloadFileFailureError
	}
	return content, nil
}

func (s service) GetDetails(mobileNumber string) (models2.QuarantineDetails, *models2.Error) {
	quarantine, err := s.repository.GetDetails(mobileNumber)
	if err != nil {
		return models2.QuarantineDetails{}, err
	}
	return utils.GetMappedQuarantine(quarantine), nil
}

func (s service) UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng string) *models2.Error {
	return s.repository.UpdateCurrentLocation(mobileNumber, currentLocationLat, currentLocationLng)
}

func (s service) UpdateDeviceTokenId(mobileNumber, deviceTokenId string) *models2.Error {
	return s.repository.UpdateDeviceTokenId(mobileNumber, deviceTokenId)
}

func (s service) NotifySO(request models2.NotificationRequest, mobileNumber string) *models2.Error {
	quarantine, err := s.repository.GetDetails(mobileNumber)
	if err != nil {
		return err
	}

	failedTokens, err := notify.SendNotification([]string{quarantine.SupervisingOfficer.DeviceTokenId}, map[string]string{
		"type":          request.Type,
		"message":       request.Message,
		"mobile_number": quarantine.MobileNumber,
		"name":          quarantine.Name,
		"address":       quarantine.Address.AddressLine1,
	})
	if err != nil {
		return err
	}
	if len(failedTokens) > 0 {
		return &constants.SendNotificationFailedError
	}
	return nil
}

func mapToDBQuarantine(detailRequest models2.QuarantineDetails) (dbmodels.Quarantine, *models2.Error) {
	DOB, err := time.Parse(constants.DetailsTimeFormat, detailRequest.DOB)
	if err != nil {
		logrus.Error("Could not parse dob ", err)
		return dbmodels.Quarantine{}, &constants.DOBIncorrectFormatError
	}

	QuarantineStartedFrom, err := time.Parse(constants.DetailsTimeFormat, detailRequest.QuarantineStartedFrom)
	if err != nil {
		logrus.Error("Could not parse quarantine started from ", err)
		return dbmodels.Quarantine{}, &constants.QuarantineDateIncorrectFormatError
	}
	history, mapErr := mapToDBTravelHistory(detailRequest.TravelHistory)
	if mapErr != nil {
		return dbmodels.Quarantine{}, mapErr
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


func mapToDBTravelHistory(travelHistoryRequest []models2.TravelHistory) ([]dbmodels.QuarantineTravelHistory, *models2.Error) {
	var travelHistory []dbmodels.QuarantineTravelHistory
	for _, history := range travelHistoryRequest {
		visitedDate, err := time.Parse(constants.DetailsTimeFormat, history.VisitDate)
		if err != nil {
			logrus.Error("Could not parse visited date of travel ", history.PlaceVisited, " error-", err)
			return nil, &constants.TravelDateIncorrectFormatError
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

func mapToDBAddress(address *models2.Address) dbmodels.QuarantineAddress {
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
