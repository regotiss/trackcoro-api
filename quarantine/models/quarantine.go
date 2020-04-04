package models

import "github.com/jinzhu/gorm"

type Quarantine struct {
	gorm.Model
	Name         string
	MobileNumber string `gorm:"unique"`
	Address      QuarantineAddress
}
