package websocket

import (
	"encoding/json"
	"log"
	"websocket-chat/internal/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan interface{}
	RoomID string
	User   *models.User
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.UnregisterClient(c)
		c.Conn.Close()
	}()
	for {
		_, messageData, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var baseMsg BaseMessage
		err = json.Unmarshal(messageData, &baseMsg)
		if err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		switch baseMsg.Type {
		case "add_option":
			var addOptionMsg AddOptionMessage
			err = json.Unmarshal(messageData, &addOptionMsg)
			if err != nil {
				log.Printf("Invalid add_option message: %v", err)
				continue
			}
			c.handleAddOption(hub, addOptionMsg)
		case "vote":
			var voteMsg VoteMessage
			err = json.Unmarshal(messageData, &voteMsg)
			if err != nil {
				log.Printf("Invalid vote message: %v", err)
				continue
			}
			c.handleVote(hub, voteMsg)
		default:
			log.Printf("Unknown message type: %s", baseMsg.Type)
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		err := c.Conn.WriteJSON(message)
		if err != nil {
			break
		}
	}
}

func (c *Client) handleAddOption(hub *Hub, msg AddOptionMessage) {
	if msg.Content == "" {
		errorMsg := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "error",
			Message: "Content cannot be empty",
		}
		c.Send <- errorMsg
		return
	}

	err := hub.SqlStore.ChangeOption(c.User.UserID, c.RoomID, msg.Content)
	if err != nil {
		errorMsg := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "error",
			Message: "Failed to create option",
		}
		c.Send <- errorMsg
		return
	}

	fullRoomStateMsg, err := hub.SqlStore.GetFullRoomState(c.RoomID)
	if err != nil {
		errorMsg := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "error",
			Message: "Failed to get room state",
		}
		c.Send <- errorMsg
		return
	}

	hub.Broadcast <- BroadcastMessage{
		RoomID:  c.RoomID,
		Message: *fullRoomStateMsg,
	}
}

func (c *Client) handleVote(hub *Hub, msg VoteMessage) {
	if msg.OptionID == "" {
		errorMsg := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "error",
			Message: "Content cannot be empty",
		}
		c.Send <- errorMsg
		return
	}

	err := hub.SqlStore.ChangeVote(c.User.UserID, msg.OptionID)
	if err != nil {
		errorMsg := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "error",
			Message: "Failed to create vote",
		}
		c.Send <- errorMsg
		return
	}

	fullRoomStateMsg, err := hub.SqlStore.GetFullRoomState(c.RoomID)
	if err != nil {
		errorMsg := struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		}{
			Type:    "error",
			Message: "Failed to get room state",
		}
		c.Send <- errorMsg
		return
	}

	hub.Broadcast <- BroadcastMessage{
		RoomID:  c.RoomID,
		Message: *fullRoomStateMsg,
	}
}
