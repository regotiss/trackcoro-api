package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trackcoro/constants"
	"trackcoro/database/models"
	models2 "trackcoro/models"
	models1 "trackcoro/quarantine/models"
	"trackcoro/token"
)

func AddTokenInHeader(ctx *gin.Context, mobileNumber string, role string) {
	tokenBody := token.UserInfo{MobileNumber: mobileNumber, Role: role}
	generatedToken, generatedTime, err := token.GenerateToken(tokenBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Header("Token", generatedToken)
	ctx.Header("Generated-At", generatedTime.String())
}

func GetMobileNumber(ctx *gin.Context) string {
	mobileNumber, _ := ctx.Get(constants.MobileNumber)
	return mobileNumber.(string)
}

func GetMappedQuarantine(quarantine models.Quarantine) models2.QuarantineDetails {
	return models2.QuarantineDetails{
		MobileNumber:           quarantine.MobileNumber,
		Name:                   quarantine.Name,
		Occupation:             quarantine.Occupation,
		DOB:                    quarantine.DOB.String(),
		Address:                mapFromDBAddress(quarantine.Address),
		TravelHistory:          mapFromDBTravelHistory(quarantine.TravelHistory),
		AnyPractitionerConsult: quarantine.AnyPractitionerConsult,
		NoOfQuarantineDays:     quarantine.NoOfQuarantineDays,
		QuarantineStartedFrom:  quarantine.QuarantineStartedFrom.String(),
		FamilyMembers:          quarantine.FamilyMembers,
		SecondaryContactNumber: quarantine.SecondaryContactNumber,
	}
}

func mapFromDBAddress(address models.QuarantineAddress) models2.Address {
	return models2.Address{
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		AddressLine3: address.AddressLine3,
		Locality:     address.Locality,
		City:         address.City,
		District:     address.District,
		State:        address.State,
		Country:      address.Country,
		PinCode:      address.PinCode,
		Coordinates: models2.Coordinates{
			Latitude:  address.Latitude,
			Longitude: address.Longitude,
		},
	}
}

func mapFromDBTravelHistory(quarantineTravelHistory []models.QuarantineTravelHistory) []models1.TravelHistory {
	var travelHistory []models1.TravelHistory
	for _, history := range quarantineTravelHistory {
		travelHistory = append(travelHistory, models1.TravelHistory{
			PlaceVisited:         history.PlaceVisited,
			VisitDate:            history.VisitDate.String(),
			TimeSpentInDays:      history.TimeSpentInDays,
			ModeOfTransportation: history.ModeOfTransportation,
		})
	}
	return travelHistory
}
