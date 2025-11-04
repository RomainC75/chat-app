package chat_client

type IWebSocket interface {
	WriteTextMessage(data []byte) error
	WriteCloseMessage() error
	GetChan() chan (ICommandMessageIn)
}
