package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"websocket-chat/internal/store"
	"websocket-chat/internal/utils"
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

		token, err := utils.GenerateJWT(user.UserID, user.RoomID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		room, _ := sqlStore.GetFullRoomState(req.RoomID)
		hub.Broadcast <- ws.BroadcastMessage{
			RoomID:  req.RoomID,
			Message: room,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}

func UpdateUserWithOption(hub *ws.Hub, sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserWithOptionRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID").(string)

		if req.RoomID == "" || req.DisplayName == "" {
			http.Error(w, "roomID and displayName are required", http.StatusBadRequest)
			return
		}

		err = sqlStore.ChangeUserName(userID, req.RoomID, req.DisplayName)
		if err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
		fmt.Println("created user")

		if len(req.OptionContent) > 0 {
			err = sqlStore.ChangeOption(userID, req.RoomID, req.OptionContent)
			if err != nil {
				http.Error(w, "Failed to update option", http.StatusInternalServerError)
				return
			}
		}

		room, _ := sqlStore.GetFullRoomState(req.RoomID)
		hub.Broadcast <- ws.BroadcastMessage{
			RoomID:  req.RoomID,
			Message: room,
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User and option updated successfully"))
	}
}

func CreateAvailability(sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateAvailabilityRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		userID := r.Context().Value("userID").(string)

		if req.RoomID == "" || userID == "" || req.Dates == nil {
			http.Error(w, "roomID and userID and dates are required", http.StatusBadRequest)
			return
		}

		err = sqlStore.DeleteUserDates(req.RoomID, userID)

		if err != nil {
			http.Error(w, "Error deleting user dates", http.StatusInternalServerError)
			return
		}

		for _, date := range req.Dates {
			_, err := sqlStore.CreateDate(req.RoomID, userID, date)
			if err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Availability created successfully"))
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

func GetDates(sqlStore *store.SQLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("roomID")
		results, err := sqlStore.GetDatesByRoomID(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(results)
	}
}
