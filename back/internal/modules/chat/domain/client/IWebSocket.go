package chat_client

import (
	"chat/internal/modules/chat/domain/messages"
)

type IWebSocket interface {
	WriteTextMessage(message *messages.Message) error
	WriteInfoMessage(messageType MessageOutType, content map[string]string) error
	WriteEvent(event IEvents) error
	WriteCloseMessage() error
	LinkToClient(c *Client)
}
