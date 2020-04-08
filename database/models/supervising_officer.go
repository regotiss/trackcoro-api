package models

import "github.com/jinzhu/gorm"

type SupervisingOfficer struct {
	gorm.Model
	MobileNumber         string `gorm:"unique;not null"`
	Name                 string
	BadgeId              string
	Designation          string
	PoliceStationAddress string
	DeviceTokenId        string
	AdminID              uint
	Quarantines          []Quarantine
}
