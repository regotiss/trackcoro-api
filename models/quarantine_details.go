package models

import "trackcoro/quarantine/models"

type QuarantineDetails struct {
	MobileNumber           string                 `json:"mobile_number"`
	Name                   string                 `json:"name" binding:"required"`
	Address                Address                `json:"address" binding:"required"`
	Occupation             string                 `json:"occupation" binding:"required"`
	DOB                    string                 `json:"date_of_birth" binding:"required"`
	TravelHistory          []models.TravelHistory `json:"travel_history"`
	AnyPractitionerConsult bool                   `json:"any_practitioner_consult" binding:"required"`
	NoOfQuarantineDays     uint                   `json:"number_of_quarantine_days" binding:"required"`
	QuarantineStartedFrom  string                 `json:"quarantine_started_from" binding:"required"`
	FamilyMembers          uint                   `json:"family_members" binding:"required"`
	SecondaryContactNumber string                 `json:"secondary_contact_number" binding:"required"`
	DeviceTokenId          string                 `json:"device_token_id" binding:"required"`
}
