package models

type AddSORequest struct {
	MobileNumber         string `json:"mobile_number" binding:"required,len=10"`
	Name                 string `json:"name" binding:"required"`
	BadgeId              string `json:"badge_id" binding:"required"`
	Designation          string `json:"designation" binding:"required"`
	PoliceStationAddress string `json:"police_station_address" binding:"required"`
}
