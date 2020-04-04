package models

import "github.com/jinzhu/gorm"

type Quarantine struct {
	gorm.Model
	MobileNumber string
}
