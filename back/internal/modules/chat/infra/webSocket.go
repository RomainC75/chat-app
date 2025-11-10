package chat_app_infra

import (
	chat_client "chat/internal/modules/chat/domain/client"
	"chat/internal/modules/chat/domain/messages"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RawMessageIn struct {
	MessageType int
	P           []byte
	Err         error
}

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
	readChan chan (chat_client.ICommandMessageIn)
	ctx      context.Context
	cancel   context.CancelFunc
	client   *chat_client.Client
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WebSocket, error) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	fws := &WebSocket{
		conn:   conn,
		ctx:    ctx,
		cancel: cancel,
	}
	fws.listenToNewMessages()
	return fws, nil
}

func (fws *WebSocket) LinkToClient(c *chat_client.Client) {
	fws.client = c
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
					// ! notify every clients tha this client disconnected
					fws.cancel()
					continue
				}

				msg, err := fws.HandleMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error handling the message in")
					fws.sendErrorMessage()
					continue
				}
				fws.client.ListenToMessageIn(msg)
			}
		}
	}()
}

func (fws *WebSocket) HandleMessageIn(payload []byte) (chat_client.ICommandMessageIn, error) {
	msg, err := UnMarshallMessageIn(payload)
	if err != nil {
		slog.Error("-> client : error unMarshalling the payload")
		return nil, err
	}
	switch msg.Type {
	case BROADCAST_MESSAGE:
		return NewBroadcastMessageIn(msg.Content["message"]), nil
	case ROOM_MESSAGE:
		roomId, err := uuid.Parse(msg.Content["room_id"])
		if err != nil {
			return nil, fmt.Errorf("unparsable room id")
		}
		return NewRoomMessageIn(roomId, msg.Content["message"]), nil
	case CREATE_ROOM:
		return NewCreateRoomICommandMessageIn(msg.Content["name"], msg.Content["description"]), nil
	case CONNECT_TO_ROOM:
		roomUuid, err := uuid.Parse(msg.Content["room_id"])
		if err != nil {
			return nil, fmt.Errorf("unparsable room id")
		}
		return NewConnectToRoomIn(roomUuid), nil
	}
	return nil, fmt.Errorf("unknown message type")

}

func (fws *WebSocket) GetChan() chan (chat_client.ICommandMessageIn) {
	return fws.readChan
}

func (fws *WebSocket) WriteCloseMessage() error {
	return fws.conn.WriteMessage(websocket.CloseMessage, nil)
}

func (fws *WebSocket) Cancel() {
	fws.cancel()
}

func (fws *WebSocket) sendErrorMessage() {
	badRequestMessage := BuildMessageOut(chat_client.ERROR, map[string]string{
		"message": "bad request",
	})
	bMessageOut, _ := json.Marshal(badRequestMessage)
	fws.conn.WriteMessage(websocket.TextMessage, bMessageOut)
}

func (fws *WebSocket) WriteTextMessage(message *messages.Message) error {
	data := BuildMessageOut(chat_client.HELLO, map[string]string{
		"message": "readyToCommunicate :-)",
	})
	bMessageOut, _ := json.Marshal(data)
	return fws.conn.WriteMessage(websocket.TextMessage, bMessageOut)
}

func (fws *WebSocket) WriteHelloMessage() error {
	data := BuildMessageOut(chat_client.HELLO, map[string]string{
		"message": "readyToCommunicate :-)",
	})
	bMessageOut, _ := json.Marshal(data)
	return fws.conn.WriteMessage(websocket.TextMessage, bMessageOut)
}

func (fws *WebSocket) WriteInfoMessage(messageType chat_client.MessageOutType, content map[string]string) error {
	data := BuildMessageOut(messageType, content)
	bMessageOut, _ := json.Marshal(data)
	return fws.conn.WriteMessage(websocket.TextMessage, bMessageOut)
}

func (fws *WebSocket) WriteEvent(event chat_client.IEvents) error {
	event.Execute(fws)
	return nil
}
