package websocket

// websocket requests
type BaseMessage struct {
	Type string `json:"type"`
}

type AddOptionMessage struct {
	Content string `json:"content"`
}

type VoteMessage struct {
	OptionID string `json:"optionID"`
}
