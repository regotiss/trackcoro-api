package models

type Address struct {
	AddressLine1 string      `json:"address_line1,omitempty" binding:"required"`
	AddressLine2 string      `json:"address_line2,omitempty"`
	AddressLine3 string      `json:"address_line3,omitempty"`
	Locality     string      `json:"locality,omitempty" binding:"required"`
	City         string      `json:"city,omitempty" binding:"required"`
	District     string      `json:"district,omitempty" binding:"required"`
	State        string      `json:"state,omitempty" binding:"required"`
	Country      string      `json:"country,omitempty" binding:"required"`
	PinCode      string      `json:"pincode,omitempty" binding:"required"`
	Coordinates  Coordinates `json:"coordinates,omitempty" binding:"required"`
}
