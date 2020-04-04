package models

type SaveDetailsRequest struct {
	MobileNumber           string          `json:"mobile_number" binding:"required"`
	Name                   string          `json:"name" binding:"required"`
	Address                Address         `json:"address" binding:"required"`
	Occupation             string          `json:"occupation" binding:"required"`
	DOB                    string          `json:"date_of_birth" binding:"required"`
	TravelHistory          []TravelHistory `json:"travel_history"`
	AnyPractitionerConsult bool            `json:"any_practitioner_consult" binding:"required"`
	NoOfQuarantineDays     uint            `json:"number_of_quarantine_days" binding:"required"`
	QuarantineStartedFrom  string          `json:"quarantine_started_from" binding:"required"`
	FamilyMembers          uint            `json:"family_members" binding:"required"`
	SecondaryContactNumber string          `json:"secondary_contact_number" binding:"required"`
}
