package models

import "github.com/jinzhu/gorm"

type QuarantineAddress struct {
	gorm.Model
	AddressLine1 string
	AddressLine2 string
	AddressLine3 string
	Locality     string
	City         string
	District     string
	State        string
	Country      string
	PinCode      string
	Latitude     string
	Longitude    string
	QuarantineID uint
}
