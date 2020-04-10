package models

type SODetails struct {
	MobileNumber         string `json:"mobile_number" binding:"required,len=10"`
	Name                 string `json:"name,omitempty" binding:"required"`
	BadgeId              string `json:"badge_id,omitempty" binding:"required"`
	Designation          string `json:"designation,omitempty" binding:"required"`
	PoliceStationAddress string `json:"police_station_address,omitempty" binding:"required"`
}
