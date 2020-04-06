package models

type Coordinates struct {
	Latitude  string `json:"latitude,omitempty" binding:"required"`
	Longitude string `json:"longitude,omitempty" binding:"required"`
}
