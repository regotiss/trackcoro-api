package models

import "trackcoro/models"

type NotifyQuarantine struct {
	models.NotificationRequest
	QuarantineRequest
}