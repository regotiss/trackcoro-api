package models

type SODetails struct {
	MobileNumber         string `form:"mobile_number" json:"mobile_number" binding:"required,len=10"`
	Name                 string `form:"name" json:"name,omitempty" binding:"required"`
	BadgeId              string `form:"badge_id" json:"badge_id,omitempty" binding:"required"`
	Designation          string `form:"designation" json:"designation,omitempty" binding:"required"`
	PoliceStationAddress string `form:"police_station_address" json:"police_station_address,omitempty" binding:"required"`
}
