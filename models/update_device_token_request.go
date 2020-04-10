package models

type DeviceTokeIdRequest struct {
	DeviceTokeId string `json:"device_token_id" binding:"required"`
}
