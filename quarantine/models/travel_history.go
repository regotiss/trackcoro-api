package models

type TravelHistory struct {
	PlaceVisited         string `json:"place_visited" binding:"required"`
	VisitDate            string `json:"visit_date" binding:"required"`
	TimeSpentInDays      uint   `json:"time_spent_in_days" binding:"required"`
	ModeOfTransportation string `json:"mode_of_transportation" binding:"required"`
}
