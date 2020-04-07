package models

type RemoveQuarantineRequest struct {
	MobileNumber string `json:"quarantine_mobile_number" binding:"required,len=10"`
}
