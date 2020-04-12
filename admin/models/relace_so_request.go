package models

type ReplaceSORequest struct {
	OldSOMobileNumber string `form:"old_so_mobile_number" json:"old_so_mobile_number" binding:"required,len=10"`
	NewSOMobileNumber string `form:"new_so_mobile_number" json:"new_so_mobile_number" binding:"required,len=10"`
}
