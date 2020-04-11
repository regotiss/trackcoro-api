package models

type QuarantineRequest struct {
	MobileNumber string `json:"quarantine_mobile_number" binding:"required,len=10"`
}
