package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type QuarantineTravelHistory struct {
	gorm.Model
	PlaceVisited         string
	VisitDate            time.Time
	TimeSpentInDays      uint
	ModeOfTransportation string
	QuarantineID         uint
}
