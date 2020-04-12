package models

type VerifyRequest struct {
	MobileNumber string `form:"mobile_number" json:"mobile_number" binding:"required,len=10"`
}