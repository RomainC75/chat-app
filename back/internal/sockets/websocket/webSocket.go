package websocket

import (
	socket_shared "chat/internal/sockets/shared"
	"context"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	conn     *websocket.Conn
	readChan chan (socket_shared.RawMessageIn)
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWebSocket(conn *websocket.Conn) *WebSocket {
	ctx, cancel := context.WithCancel(context.Background())

	fws := &WebSocket{
		conn:     conn,
		readChan: make(chan (socket_shared.RawMessageIn)),
		ctx:      ctx,
		cancel:   cancel,
	}
	fws.listenToNewMessages()
	return fws
}

func (fws *WebSocket) listenToNewMessages() {
	go func() {
		for {
			select {
			case <-fws.ctx.Done():
				return
			default:
				mType, payload, err := fws.conn.ReadMessage()
				messageIn := socket_shared.RawMessageIn{
					MessageType: mType,
					P:           payload,
					Err:         err,
				}
				fws.readChan <- messageIn
			}
		}
	}()
}

func (fws *WebSocket) ReadMessage() chan (socket_shared.RawMessageIn) {
	return fws.readChan
}

func (fws *WebSocket) WriteMessage(messageType int, data []byte) error {
	return fws.conn.WriteMessage(socket_shared.TextMessage, data)
}

func (fws *WebSocket) Cancel() {
	fws.cancel()
}
