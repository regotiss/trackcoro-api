package models

type GetQuarantinesRequest struct {
	MobileNumber string `json:"so_mobile_number" binding:"required,len=10"`
}
