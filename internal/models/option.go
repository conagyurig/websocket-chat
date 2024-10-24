package models

type Option struct {
	OptionID string `json:"id"`
	RoomID   string `json:"roomId"`
	UserID   string `json:"userId"`
	Content  string `json:"content"`
}
