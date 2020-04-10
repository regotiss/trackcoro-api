package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Quarantine struct {
	gorm.Model
	MobileNumber             string `gorm:"unique;not null"`
	Name                     string
	Address                  QuarantineAddress
	TravelHistory            []QuarantineTravelHistory
	Occupation               string
	DOB                      time.Time
	AnyPractitionerConsult   bool
	NoOfQuarantineDays       uint
	QuarantineStartedFrom    time.Time
	FamilyMembers            uint
	SecondaryContactNumber   string
	DeviceTokenId            string
	CurrentLocationLatitude  string
	CurrentLocationLongitude string
	SupervisingOfficer       *SupervisingOfficer
	SupervisingOfficerID     uint
	PhotoUploads             []*PhotoUpload
}
