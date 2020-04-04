package quarantine

import "trackcoro/quarantine/models"

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
	user := models.Quarantine{Name: detailsRequest.Name, MobileNumber: detailsRequest.MobileNumber}
	return s.repository.SaveDetails(user)
}

func NewService(repository Repository) Service {
	return service{repository}
}
