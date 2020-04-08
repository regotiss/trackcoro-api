package models

type NotificationRequest struct {
	Type string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}