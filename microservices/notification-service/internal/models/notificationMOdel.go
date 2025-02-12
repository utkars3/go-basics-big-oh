package models

type Notification struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}
