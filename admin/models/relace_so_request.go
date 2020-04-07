package models

type ReplaceSORequest struct {
	OldSOMobileNumber string `json:"old_so_mobile_number" binding:"required,len=10"`
	NewSOMobileNumber string `json:"new_so_mobile_number" binding:"required,len=10"`
}
