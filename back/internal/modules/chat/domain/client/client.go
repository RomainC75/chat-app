package client

import (
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/google/uuid"
)

type IManager interface {
	RemoveClient(*Client)
	SendBroadcastMessage(userData socket_shared.UserData, msgIn messages.MessageIn)
	SendRoomMessage(c *Client, roomId string, message string)
	CreateRoom(c *Client, roomName string)
	ConnectUserAndRoom(c *Client, roomId uuid.UUID) error
}

type IRoom interface {
	GetId() uuid.UUID
	GetClients() []socket_shared.UserData
}

type Client struct {
	manager  IManager
	room     IRoom
	conn     socket_shared.IWebSocket
	user     socket_shared.UserData
	egress   chan ([]byte)
	cancelFn context.CancelFunc
	ctx      context.Context
}

func NewClient(manager IManager, conn socket_shared.IWebSocket, userData socket_shared.UserData) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		manager:  manager,
		conn:     conn,
		user:     userData,
		egress:   make(chan []byte),
		cancelFn: cancel,
		ctx:      ctx,
	}
}

func (c *Client) PrepareToBeDeleted() {
	c.cancelFn()
}

func (c *Client) SendToClient(msg messages.MessageOut) {
	m, _ := json.Marshal(msg)
	c.egress <- m
}

func (c *Client) GoListen() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			case rawMessageIn := <-c.conn.GetChan():
				err := rawMessageIn.Err
				payload := rawMessageIn.P
				if err != nil {
					slog.Error("client disconnected", "err", err)
					c.manager.RemoveClient(c)
				}
				fmt.Println("----------> Message ", string(payload))
				message, err := messages.UnMarshallMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error unMarshalling the payload")
				}
				c.HandleMessageIn(message)
			}
		}
	}()
}

func (c *Client) GoWrite() {
	c.writeHelloMessage()
	go func() {
		defer func() {
			c.manager.RemoveClient(c)
		}()

		for {
			select {
			case <-c.ctx.Done():
				return
			case message, ok := <-c.egress:
				if !ok {
					if err := c.conn.WriteMessage(socket_shared.CloseMessage, nil); err != nil {
						log.Println("connection closed:", err)
					}
					continue
				}

				if err := c.conn.WriteMessage(socket_shared.TextMessage, message); err != nil {
					log.Printf("failed to send message: %v\n", err)
				}
				log.Println("message sent")
			}
		}
	}()
}

func (c *Client) GetUserData() socket_shared.UserData {
	return c.user
}

func (c *Client) HandleMessageIn(msg messages.MessageIn) {
	switch msg.Type {
	case messages.BROADCAST_MESSAGE:
		c.manager.SendBroadcastMessage(c.user, msg)
	case messages.ROOM_MESSAGE:
		c.manager.SendRoomMessage(c, msg.Content["room_id"], msg.Content["message"])
	case messages.CREATE_ROOM:
		c.manager.CreateRoom(c, msg.Content["name"])
	case messages.CONNECT_TO_ROOM:
		roomIdStr := msg.Content["room_id"]
		roomId, _ := uuid.Parse(roomIdStr)
		_ = c.manager.ConnectUserAndRoom(c, roomId)
	default:
		c.writeErrorMessage()
		return
	}
}

func (c *Client) writeHelloMessage() {
	helloMessage := messages.BuildMessageOut(messages.HELLO, map[string]string{
		"message": "readyToCommunicate :-)",
	})
	bMessageOut, _ := json.Marshal(helloMessage)
	c.conn.WriteMessage(socket_shared.TextMessage, bMessageOut)
}

func (c *Client) writeErrorMessage() {
	badRequestMessage := messages.BuildMessageOut(messages.ERROR, map[string]string{
		"message": "bad request",
	})
	bMessageOut, _ := json.Marshal(badRequestMessage)
	c.conn.WriteMessage(socket_shared.TextMessage, bMessageOut)
}

func (c *Client) ConnectToRoom(room IRoom) {
	c.room = room
	roomUsers := room.GetClients()
	message := messages.BuildConnectedToRoomMessageOut(roomUsers, room.GetId())
	bMessageOut, _ := json.Marshal(message)
	c.conn.WriteMessage(socket_shared.TextMessage, bMessageOut)
}
