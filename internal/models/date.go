package models

type Date struct {
	DateID string `json:"id"`
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
	Date   string `json:"date"`
}
