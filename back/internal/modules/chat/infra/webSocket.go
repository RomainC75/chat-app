package chat_app_infra

import (
	"chat/internal/modules/chat/domain/messages"
	chat_socket "chat/internal/modules/chat/domain/socket"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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
	readChan chan (chat_socket.ICommandMessageIn)
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
		readChan: make(chan (chat_socket.ICommandMessageIn)),
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
				_, payload, err := fws.conn.ReadMessage()
				if err != nil {
					slog.Error("-> client : error reading message from websocket", slog.String("error", err.Error()))
					fws.sendErrorMessage()
					continue
				}

				msg, err := fws.HandleMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error handling the message in")
					fws.sendErrorMessage()
					continue
				}
				fws.readChan <- msg
			}
		}
	}()
}

func (fws *WebSocket) HandleMessageIn(payload []byte) (chat_socket.ICommandMessageIn, error) {
	msg, err := messages.UnMarshallMessageIn(payload)
	if err != nil {
		slog.Error("-> client : error unMarshalling the payload")
	}
	switch msg.Type {
	case messages.BROADCAST_MESSAGE:
		return chat_socket.NewBroadcastMessageIn(msg.Content["message"]), nil
	case messages.ROOM_MESSAGE:
		roomId, err := uuid.Parse(msg.Content["room_id"])
		if err != nil {
			return nil, fmt.Errorf("unparsable room id")
		}
		return chat_socket.NewSendRoomMessageIn(roomId, msg.Content["message"]), nil
	case messages.CREATE_ROOM:
		return chat_socket.NewCreateRoomICommandMessageIn(msg.Content["name"], msg.Content["description"]), nil
	case messages.CONNECT_TO_ROOM:
		return chat_socket.NewConnectToRoomIn(msg.Content["room_id"]), nil
	}
	return nil, fmt.Errorf("unknown message type")

}

func (fws *WebSocket) GetChan() chan (chat_socket.ICommandMessageIn) {
	return fws.readChan
}

func (fws *WebSocket) WriteMessage(messageType int, data []byte) error {
	fmt.Println("-------> message to SEND BACK: ", string(data))
	return fws.conn.WriteMessage(chat_socket.TextMessage, data)
}

func (fws *WebSocket) Cancel() {
	fws.cancel()
}

func (fws *WebSocket) sendErrorMessage() {
	badRequestMessage := chat_socket.BuildMessageOut(messages.ERROR, map[string]string{
		"message": "bad request",
	})
	bMessageOut, _ := json.Marshal(badRequestMessage)
	fws.conn.WriteMessage(chat_socket.TextMessage, bMessageOut)
}
