package quarantine

type Service interface {
	Verify(mobileNumber string) (bool, error)
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) (bool, error) {
	return s.repository.isExists(mobileNumber)
}

func NewService(repository Repository) Service {
	return service{repository}
}