package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	socket_shared "sockets/shared"
)

type IManager interface {
	SendBroadcastMessage(userData socket_shared.UserData, msgIn MessageIn)
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
	rawMessageInChan := c.conn.ReadMessage()
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			// ! chan
			case rawMessageIn := <-rawMessageInChan:
				payload := rawMessageIn.P
				err := rawMessageIn.Err
				if err != nil {
					slog.Error("client disconnected %s", err.Error())

				}
				fmt.Println("----------> Message ", string(payload))
				message, err := UnMarshallMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error unMarshalling the payload")
				}
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

		for {
			message, ok := <-c.egress
			if !ok {
				if err := c.conn.WriteMessage(socket_shared.CloseMessage, nil); err != nil {
					log.Println("connection closed:", err)
				}
				break
			}

			if err := c.conn.WriteMessage(socket_shared.TextMessage, message); err != nil {
				log.Printf("failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}()
}

func (c *Client) HandleMessageIn(msg MessageIn) {
	fmt.Println("msg", msg.Type, msg.Content["name"])
	switch msg.Type {
	case BROADCAST_MESSAGE:
		c.manager.SendBroadcastMessage(c.user, msg)
		fmt.Println("BROADCAST : ", msg)
	case ROOM_MESSAGE:
		fmt.Println("message room")
	default:
		return
	}
	// ! Release
	c.conn.MessageInTreated()
}

func (c *Client) GetUserData() socket_shared.UserData {
	return c.user
}
