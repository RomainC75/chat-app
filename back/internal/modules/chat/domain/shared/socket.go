package socket_shared

type RawMessageIn struct {
	MessageType int
	P           []byte
	Err         error
}

type IWebSocket interface {
	WriteMessage(messageType int, data []byte) error
	GetChan() chan (RawMessageIn)
}

// copied from gorilla/websocket to purify the socket domain
const (
	TextMessage   = 1
	BinaryMessage = 2
	CloseMessage  = 8
	PingMessage   = 9
	PongMessage   = 10
)
