package handlers

import (
	"encoding/json"
	"net/http"

	"websocket-chat/internal/store"
	ws "websocket-chat/internal/websocket"
)

func CreateRoom(sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRoomRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.RoomName == "" {
			http.Error(w, "roomName is required", http.StatusBadRequest)
			return
		}

		room, err := sqlStore.CreateRoom(req.RoomName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(room)
	}
}

func CreateUser(sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.RoomID == "" || req.DisplayName == "" {
			http.Error(w, "roomID and displayName are required", http.StatusBadRequest)
			return
		}

		user, err := sqlStore.CreateUser(req.RoomID, req.DisplayName)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUserWithOption(hub *ws.Hub, sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserWithOptionRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.RoomID == "" || req.DisplayName == "" {
			http.Error(w, "roomID and displayName are required", http.StatusBadRequest)
			return
		}

		user, err := sqlStore.CreateUser(req.RoomID, req.DisplayName)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		if len(req.OptionContent) > 0 {
			_, err = sqlStore.CreateOption(req.RoomID, user.UserID, req.OptionContent)
			if err != nil {
				http.Error(w, "Failed to create option", http.StatusInternalServerError)
				return
			}
		}

		room, _ := sqlStore.GetFullRoomState(req.RoomID)
		hub.Broadcast <- ws.BroadcastMessage{
			RoomID:  req.RoomID,
			Message: room,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func GetRoomState(sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("roomID")
		room, err := sqlStore.GetFullRoomState(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(room)
	}
}
