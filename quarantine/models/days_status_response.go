package models

import "time"

type DaysStatusResponse struct {
	NoOfQuarantineDays uint      `json:"no_of_quarantine_days"`
	RemainingDays      int      `json:"remaining_days"`
	StatedFrom         time.Time `json:"started_from"`
}
