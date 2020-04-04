package quarantine

import (
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
	user := models.Quarantine{Name: detailsRequest.Name, MobileNumber: detailsRequest.MobileNumber,
		Address: mapAddress(detailsRequest.Address)}
	return s.repository.SaveDetails(user)
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
