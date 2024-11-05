package handlers

import (
	"net/http"
	"websocket-chat/internal/utils"
	ws "websocket-chat/internal/websocket"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *ws.Hub, w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomID")
	tokenString := r.URL.Query().Get("token")

	if tokenString == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}

	claims := &utils.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return utils.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID

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
