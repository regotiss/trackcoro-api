package models

import "trackcoro/models"

type QVerifyResponse struct {
	models.VerifyResponse
	IsSignUpCompleted bool `json:"is_sign_up_completed"`
}
