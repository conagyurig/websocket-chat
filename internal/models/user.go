package models

type User struct {
	UserID      string `json:"id"`
	RoomID      string `json:"roomId"`
	DisplayName string `json:"name"`
}
