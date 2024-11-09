package store

import "websocket-chat/internal/models"

type FullRoomStateMessage struct {
	RoomName    string          `json:"roomName"`
	Users       []models.User   `json:"users"`
	Options     []models.Option `json:"options"`
	Votes       []models.Vote   `json:"votes"`
	RevealVotes bool            `json:"revealVotes"`
}
