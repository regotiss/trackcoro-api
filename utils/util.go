package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"trackcoro/config"
	"trackcoro/constants"
	"trackcoro/database/models"
	models2 "trackcoro/models"
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

func GetMappedQuarantines(quarantines []models.Quarantine) []models2.QuarantineDetails {
	var quarantineDetails []models2.QuarantineDetails
	for _, quarantine := range quarantines {
		quarantineDetails = append(quarantineDetails, getQuarantineDetails(quarantine))
	}
	if quarantineDetails == nil {
		quarantineDetails = []models2.QuarantineDetails{}
	}
	return quarantineDetails
}

func getQuarantineDetails(quarantine models.Quarantine) models2.QuarantineDetails {
	mappedQuarantine := GetMappedQuarantine(quarantine)
	mappedQuarantine.CurrentLocation = MapCoordinates(quarantine.CurrentLocationLatitude, quarantine.CurrentLocationLongitude)
	return mappedQuarantine
}

func GetMappedQuarantine(quarantine models.Quarantine) models2.QuarantineDetails {
	var dob string
	if !quarantine.DOB.IsZero() {
		dob = quarantine.DOB.String()
	}
	var quarantineStartedFrom string
	if !quarantine.QuarantineStartedFrom.IsZero() {
		quarantineStartedFrom = quarantine.QuarantineStartedFrom.String()
	}
	var soDetails *models2.SODetails
	if quarantine.SupervisingOfficer != nil {
		soDetails = &models2.SODetails{MobileNumber: quarantine.SupervisingOfficer.MobileNumber,
			Name: quarantine.SupervisingOfficer.Name}
	}
	photoURL := fmt.Sprintf("%s/%s", config.Config.FileServerURL, quarantine.MobileNumber)
	return models2.QuarantineDetails{
		MobileNumber:           quarantine.MobileNumber,
		Name:                   quarantine.Name,
		Occupation:             quarantine.Occupation,
		DOB:                    dob,
		Address:                mapFromDBAddress(quarantine.Address),
		TravelHistory:          mapFromDBTravelHistory(quarantine.TravelHistory),
		AnyPractitionerConsult: quarantine.AnyPractitionerConsult,
		NoOfQuarantineDays:     quarantine.NoOfQuarantineDays,
		QuarantineStartedFrom:  quarantineStartedFrom,
		FamilyMembers:          quarantine.FamilyMembers,
		SecondaryContactNumber: quarantine.SecondaryContactNumber,
		SODetails:              soDetails,
		PhotoURL:               photoURL,
	}
}

func mapFromDBAddress(address models.QuarantineAddress) *models2.Address {
	if address.ID == 0 {
		return nil
	}
	return &models2.Address{
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		AddressLine3: address.AddressLine3,
		Locality:     address.Locality,
		City:         address.City,
		District:     address.District,
		State:        address.State,
		Country:      address.Country,
		PinCode:      address.PinCode,
		Coordinates:  MapCoordinates(address.Latitude, address.Longitude),
	}
}

func MapCoordinates(latitude, longitude string) *models2.Coordinates {
	if latitude != constants.Empty || longitude != constants.Empty {
		return &models2.Coordinates{
			Latitude:  latitude,
			Longitude: longitude,
		}
	}
	return nil
}

func mapFromDBTravelHistory(quarantineTravelHistory []models.QuarantineTravelHistory) []models2.TravelHistory {
	var travelHistory []models2.TravelHistory
	for _, history := range quarantineTravelHistory {
		travelHistory = append(travelHistory, models2.TravelHistory{
			PlaceVisited:         history.PlaceVisited,
			VisitDate:            history.VisitDate.String(),
			TimeSpentInDays:      history.TimeSpentInDays,
			ModeOfTransportation: history.ModeOfTransportation,
		})
	}
	return travelHistory
}

func HandleResponse(ctx *gin.Context, err *models2.Error, response interface{}, errorHandler func(error2 *models2.Error) int) {
	code := errorHandler(err)
	if code != http.StatusOK {
		ctx.AbortWithStatusJSON(code, err)
		return
	}
	if response != nil {
		ctx.JSON(code, response)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}
