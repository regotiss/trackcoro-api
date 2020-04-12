package models

type GetQuarantinesRequest struct {
	MobileNumber string `form:"so_mobile_number" json:"so_mobile_number" binding:"required,len=10"`
}
