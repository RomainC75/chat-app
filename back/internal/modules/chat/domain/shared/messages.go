package chat_shared

type Message struct {
	Message string `json:"message"`
}

type WebSocketMessage struct {
	Type    string            `json:"type"`
	Content map[string]string `json:"content"`
}
