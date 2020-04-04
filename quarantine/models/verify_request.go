package models

type VerifyRequest struct {
	MobileNumber string `json:"mobile_number" binding:"required,len=10"`
}