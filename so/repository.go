package so

import (
	"github.com/jinzhu/gorm"
	"trackcoro/database/models"
	"trackcoro/utils"
)

type Repository interface {
	IsExists(mobileNumber string) bool
	AddQuarantine(mobileNumber string, quarantine models.Quarantine) error
	GetQuarantines(mobileNumber string) ([]models.Quarantine, error)
}

type repository struct {
	db *gorm.DB
}

func (r repository) IsExists(mobileNumber string) bool {
	user, err := utils.GetSOBy(r.db, mobileNumber)
	if err != nil {
		return false
	}
	return user.MobileNumber == mobileNumber
}

func (r repository) AddQuarantine(mobileNumber string, quarantine models.Quarantine) error {
	existingSO, err := utils.GetSOBy(r.db, mobileNumber)
	if err != nil {
		return err
	}
	quarantine.SupervisingOfficerID = existingSO.ID
	return r.db.Save(&quarantine).Error
}

func (r repository) GetQuarantines(mobileNumber string) ([]models.Quarantine, error) {
	return utils.GetQuarantines(r.db, mobileNumber)
}

func NewRepository(db *gorm.DB) Repository {
	return repository{db}
}
