package so

type Service interface {
	Verify(mobileNumber string) bool
}

type service struct {
	repository Repository
}

func (s service) Verify(mobileNumber string) bool {
	return s.repository.IsExists(mobileNumber)
}

func NewService(repository Repository) Service {
	return service{repository}
}

