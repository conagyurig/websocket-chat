package models

type DateWithUsers struct {
	Date  string `json:"date"`
	Users []User `json:"users"`
}

type RoomDatesResponse struct {
	RoomID string          `json:"roomId"`
	Dates  []DateWithUsers `json:"dates"`
}
