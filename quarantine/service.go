package quarantine

import (
	"errors"
	"github.com/sirupsen/logrus"
	"time"
	"trackcoro/quarantine/models"
)

type Service interface {
	Verify(mobileNumber string) (bool, error)
	SaveDetails(request models.SaveDetailsRequest) error
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) (bool, error) {
	return s.repository.isExists(mobileNumber)
}

func (s service) SaveDetails(detailsRequest models.SaveDetailsRequest) error {
	user, err := mapQuarantine(detailsRequest)
	if err != nil {
		return err
	}
	return s.repository.SaveDetails(user)
}

func mapQuarantine(detailRequest models.SaveDetailsRequest) (models.Quarantine, error) {
	DOB, err := time.Parse(DetailsTimeFormat, detailRequest.DOB)
	if err != nil {
		logrus.Error("Could not parse dob time", err)
		return models.Quarantine{}, errors.New(TimeParseError)
	}
	QuarantineStartedFrom, err := time.Parse(DetailsTimeFormat, detailRequest.QuarantineStartedFrom)

	if err != nil {
		logrus.Error("Could not parse quarantine started from time", err)
		return models.Quarantine{}, errors.New(TimeParseError)
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
	}, nil
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
