package models

type Vote struct {
	VoteID   string `json:"id"`
	OptionID string `json:"optionId"`
	UserID   string `json:"userId"`
}
