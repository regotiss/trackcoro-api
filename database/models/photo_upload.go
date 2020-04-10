package models

import "time"

type PhotoUpload struct {
	ID           uint64
	RequestedOn  time.Time `gorm:"not null"`
	UploadedOn   time.Time `gorm:"not null"`
	QuarantineID uint
}
