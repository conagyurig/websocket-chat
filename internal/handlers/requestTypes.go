package handlers

// http requests
type CreateUserRequest struct {
	RoomID      string `json:"roomID"`
	DisplayName string `json:"displayName"`
}

type CreateUserWithOptionRequest struct {
	RoomID        string `json:"roomID"`
	DisplayName   string `json:"displayName"`
	OptionContent string `json:"optionContent"`
}

type CreateRoomRequest struct {
	RoomName string `json:"roomName"`
}
