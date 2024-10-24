package handlers

import (
	"net/http"
	ws "websocket-chat/internal/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomID")
	userID := r.URL.Query().Get("userID")

	room, err := hub.SqlStore.GetRoomByID(roomID)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	user, err := hub.SqlStore.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	client := &ws.Client{
		Conn:   conn,
		Send:   make(chan interface{}),
		RoomID: room.RoomID,
		User:   user,
	}
	hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump(hub)
}
