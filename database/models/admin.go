package models

import "github.com/jinzhu/gorm"

type Admin struct {
	gorm.Model
	MobileNumber string `gorm:"unique;not null"`
	SupervisingOfficers []SupervisingOfficer
}
