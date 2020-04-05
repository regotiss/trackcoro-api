package models

import "github.com/jinzhu/gorm"

type SupervisingOfficer struct {
	gorm.Model
	MobileNumber         string `gorm:"unique"`
	Name                 string
	BadgeId              string
	Designation          string
	PoliceStationAddress string
	AdminID              uint
	Quarantines          []Quarantine
}
