package chat_infra

import (
	"chat/internal/modules/chat/domain/client"
	socket_shared "chat/internal/modules/chat/domain/shared"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func UnMarshallMessageIn(payload []byte) (MessageIn, error) {
	msi := MessageIn{}
	err := json.Unmarshal(payload, &msi)
	if err != nil {
		return MessageIn{}, err
	}
	return msi, err
}

// IN
type MessageIn struct {
	Type    MessageInType     `json:"type"`
	Content map[string]string `json:"content"`
}

type MessageInType string

const (
	ROOM_MESSAGE         MessageInType = "ROOM_MESSAGE"
	BROADCAST_MESSAGE    MessageInType = "BROADCAST_MESSAGE"
	CONNECT_TO_ROOM      MessageInType = "CONNECT_TO_ROOM"
	CREATE_ROOM          MessageInType = "CREATE_ROOM"
	SEND_TO_ROOM         MessageInType = "SEND_TO_ROOM"
	DISCONNECT_FROM_ROOM MessageInType = "DISCONNECT_FROM_ROOM"
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

// =======================

type WebSocket struct {
	conn     *websocket.Conn
	readChan chan (client.CommandMessageIn)
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
		readChan: make(chan (client.CommandMessageIn)),
		ctx:      ctx,
		cancel:   cancel,
	}
	fws.listenToNewMessages()
	return fws, nil
}

// ======================

func (fws *WebSocket) listenToNewMessages() {
	go func() {
		for {
			select {
			case <-fws.ctx.Done():
				return
			default:
				_, payload, err := fws.conn.ReadMessage()
				message, err := UnMarshallMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error unMarshalling the payload")
				}
				commandMessage, err := fws.handleMessageIn(message)
				if err != nil {
					slog.Error("-> client : error handling the message")
					fws.writeErrorMessage()
					continue
				}
				fws.readChan <- commandMessage

			}
		}
	}()
}

func (fws *WebSocket) handleMessageIn(msg MessageIn) (client.CommandMessageIn, error) {
	switch msg.Type {
	case BROADCAST_MESSAGE:
		return NewBroadcastMessageIn(msg.Content["message"]), nil
	case ROOM_MESSAGE:
		// c.manager.SendRoomMessage(c, msg.Content["room_id"], msg.Content["message"])
		return NewRoomMessageIn(msg.Content["room_id"], msg.Content["message"]), nil
	case CREATE_ROOM:
		// c.manager.CreateRoom(c, msg.Content["name"])
	case CONNECT_TO_ROOM:
		roomIdStr := msg.Content["room_id"]
		roomId, _ := uuid.Parse(roomIdStr)
		// _ = c.manager.ConnectUserAndRoom(c, roomId)
		return NewConnectToRoomMessageIn(roomId), nil
	}
	return nil, fmt.Errorf("invalid message type")
}

func (fws *WebSocket) GetChan() chan (client.CommandMessageIn) {
	return fws.readChan
}

func (fws *WebSocket) WriteMessage(messageType int, data []byte) error {
	fmt.Println("-------> message to SEND BACK: ", string(data))
	return fws.conn.WriteMessage(socket_shared.TextMessage, data)
}

func (fws *WebSocket) Cancel() {
	fws.cancel()
}

func (fws *WebSocket) writeErrorMessage() {
	_ = fws.conn.WriteMessage(socket_shared.TextMessage, []byte(`{"type":"ERROR","content":{"message":"Invalid message"}}`))
}
