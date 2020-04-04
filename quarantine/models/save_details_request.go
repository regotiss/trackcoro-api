package models

type SaveDetailsRequest struct {
	MobileNumber string `json:"mobile_number" binding:"required"`
	Name         string `json:"name" binding:"required"`
}
