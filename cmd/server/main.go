package main

import (
	"log"
	"net/http"
	"websocket-chat/internal/handlers"
	"websocket-chat/internal/middleware"
	"websocket-chat/internal/utils"
	"websocket-chat/internal/websocket"

	"github.com/gorilla/mux"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {

	sqlStore, err := utils.InitialiseDb()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer sqlStore.DB.Close()

	hub := websocket.NewHub(sqlStore)
	go hub.Run()

	router := mux.NewRouter()

	router.HandleFunc("/userOption", handlers.CreateUserWithOption(hub, sqlStore)).Methods("POST")
	router.HandleFunc("/rooms", handlers.CreateRoom(sqlStore)).Methods("POST")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWS(hub, w, r)
	})

	protected := router.PathPrefix("/").Subrouter()
	protected.Use(middleware.JWTAuthMiddleware)
	protected.HandleFunc("/userOption", handlers.UpdateUserWithOption(hub, sqlStore)).Methods("PUT")
	protected.HandleFunc("/userAvailability", handlers.CreateAvailability(sqlStore)).Methods("POST")
	protected.HandleFunc("/roomState", handlers.GetRoomState(sqlStore)).Methods("GET")
	protected.HandleFunc("/dates", handlers.GetDates(sqlStore)).Methods("GET")

	corsRouter := enableCORS(router)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", corsRouter)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
