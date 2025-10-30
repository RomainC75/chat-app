package websocket

import (
	socket_shared "chat/internal/sockets/shared"
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// origin := r.Header.Get("Origin")
			// cfg := config.Get()
			// frontUrl := cfg.Front.Host
			// return origin == frontUrl
			return true
		},
	}
)

type WebSocket struct {
	conn     *websocket.Conn
	readChan chan (socket_shared.RawMessageIn)
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	fws := &WebSocket{
		conn:     conn,
		readChan: make(chan (socket_shared.RawMessageIn)),
		ctx:      ctx,
		cancel:   cancel,
	}
	fws.listenToNewMessages()
	return fws, nil
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

func (fws *WebSocket) GetChan() chan (socket_shared.RawMessageIn) {
	return fws.readChan
}

func (fws *WebSocket) WriteMessage(messageType int, data []byte) error {
	fmt.Println("-------> message to SEND BACK: ", string(data))
	return fws.conn.WriteMessage(socket_shared.TextMessage, data)
}

func (fws *WebSocket) Cancel() {
	fws.cancel()
}
