package websocket

import "websocket-chat/internal/store"

type BroadcastMessage struct {
	RoomID  string
	Message interface{}
}

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
	SqlStore   *store.SQLStore
}

func NewHub(sqlStore *store.SQLStore) *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMessage),
		SqlStore:   sqlStore,
	}
}

func (h *Hub) RegisterClient(client *Client) {
	h.Register <- client
}

func (h *Hub) UnregisterClient(client *Client) {
	h.Unregister <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case broadcast := <-h.Broadcast:
			for client := range h.Clients {
				if client.RoomID == broadcast.RoomID {
					select {
					case client.Send <- broadcast.Message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
