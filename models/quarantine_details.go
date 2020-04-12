package models

type QuarantineDetails struct {
	MobileNumber           string          `json:"mobile_number"`
	Name                   string          `json:"name,omitempty" binding:"required"`
	Address                *Address        `json:"address,omitempty" binding:"required"`
	Occupation             string          `json:"occupation,omitempty" binding:"required"`
	DOB                    string          `json:"date_of_birth,omitempty" binding:"required"`
	TravelHistory          []TravelHistory `json:"travel_history,omitempty"`
	AnyPractitionerConsult bool            `json:"any_practitioner_consult,omitempty" binding:"required"`
	NoOfQuarantineDays     uint            `json:"number_of_quarantine_days,omitempty" binding:"required"`
	QuarantineStartedFrom  string          `json:"quarantine_started_from,omitempty" binding:"required"`
	FamilyMembers          uint            `json:"family_members,omitempty" binding:"required"`
	SecondaryContactNumber string          `json:"secondary_contact_number,omitempty" binding:"required"`
	DeviceTokenId          string          `json:"device_token_id,omitempty" binding:"required"`
	PhotoURL               string          `json:"photo_url,omitempty"`
	CurrentLocation        *Coordinates    `json:"current_location,omitempty"`
	SODetails              *SODetails      `json:"supervising_officer,omitempty"`
	IsPhotoUploaded        bool            `json:"isPhotoUploaded"`
}
