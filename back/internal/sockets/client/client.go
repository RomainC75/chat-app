package client

import (
	socket_shared "chat/internal/sockets/shared"
	"chat/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
)

type IManager interface {
	RemoveClient(*Client)
	SendBroadcastMessage(userData socket_shared.UserData, msgIn MessageIn)
	SendRoomMessage(msgIn MessageIn)
	CreateRoom(c *Client, roomName string)
}

type Client struct {
	manager  IManager
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

func (c *Client) SendToClient(msg MessageOut) {
	m, _ := json.Marshal(msg)
	c.egress <- m
}

func (c *Client) GoListen() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				_, payload, err := c.conn.ReadMessage()
				if err != nil {
					slog.Error("client disconnected", err.Error())
					c.manager.RemoveClient(c)

				}
				fmt.Println("----------> Message ", string(payload))
				message, err := UnMarshallMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error unMarshalling the payload")
				}
				utils.PrettyDisplay("message in : ", message)
				c.HandleMessageIn(message)
			}
		}
	}()
}

func (c *Client) GoWrite() {
	helloMessage := CreateMessageOut(HELLO, map[string]string{
		"message": "readyToCommunicate :-)",
	})
	bMessageOut, _ := json.Marshal(helloMessage)
	c.conn.WriteMessage(socket_shared.TextMessage, bMessageOut)
	go func() {
		defer func() {
			c.manager.RemoveClient(c)
		}()

		for {
			message, ok := <-c.egress
			if !ok {
				if err := c.conn.WriteMessage(socket_shared.CloseMessage, nil); err != nil {
					log.Println("connection closed:", err)
				}
				break
			}

			if err := c.conn.WriteMessage(socket_shared.TextMessage, message); err != nil {
				log.Println("failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}()
}

func (c *Client) HandleMessageIn(msg MessageIn) {
	fmt.Println("msg", msg.Type, msg.Content["name"])
	switch msg.Type {
	case BROADCAST_MESSAGE:
		fmt.Println("BROADCAST : ", msg)
		c.manager.SendBroadcastMessage(c.user, msg)
	case ROOM_MESSAGE:
		c.manager.SendRoomMessage(msg)
	case CREATE_ROOM:
		c.manager.CreateRoom(c, msg.Content["name"])
	default:
		return
	}
}
