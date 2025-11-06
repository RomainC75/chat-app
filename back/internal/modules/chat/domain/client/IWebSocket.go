package chat_client

import "chat/internal/modules/chat/domain/messages"

type IWebSocket interface {
	WriteTextMessage(message *messages.Message) error
	WriteInfoMessage(messageType string, content map[string]string) error
	WriteEvent(event IEvents) error
	WriteHelloMessage() error
	WriteCloseMessage() error
	GetChan() chan (ICommandMessageIn)
}
