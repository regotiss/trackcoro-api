package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Quarantine struct {
	gorm.Model
	MobileNumber           string `gorm:"unique"`
	Name                   string
	Address                QuarantineAddress
	Occupation             string
	DOB                    time.Time
	AnyPractitionerConsult bool
	NoOfQuarantineDays     uint
	QuarantineStartedFrom  time.Time
	FamilyMembers          uint
	SecondaryContactNumber string
}
