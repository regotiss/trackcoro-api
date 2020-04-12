package models

type QuarantineRequest struct {
	MobileNumber string `form:"quarantine_mobile_number" json:"quarantine_mobile_number" binding:"required,len=10"`
}
