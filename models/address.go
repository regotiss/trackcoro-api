package models

type Address struct {
	AddressLine1 string      `json:"address_line1" binding:"required"`
	AddressLine2 string      `json:"address_line2"`
	AddressLine3 string      `json:"address_line3"`
	Locality     string      `json:"locality" binding:"required"`
	City         string      `json:"city" binding:"required"`
	District     string      `json:"district" binding:"required"`
	State        string      `json:"state" binding:"required"`
	Country      string      `json:"country" binding:"required"`
	PinCode      string      `json:"pincode" binding:"required"`
	Coordinates  Coordinates `json:"coordinates" binding:"required"`
}
