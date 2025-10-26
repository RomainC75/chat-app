package client

import (
	socket_shared "chat/internal/sockets/shared"
	"chat/utils"
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/gorilla/websocket"
)

type IManager interface {
	RemoveClient(*Client)
}

type Client struct {
	manager  IManager
	conn     *websocket.Conn
	user     socket_shared.UserData
	egress   chan ([]byte)
	cancelFn context.CancelFunc
	ctx      context.Context
}

func NewClient(manager IManager, conn *websocket.Conn, userData socket_shared.UserData) *Client {
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

func (c *Client) KillGoroutineChildren() {
	c.cancelFn()
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
				fmt.Println("-> Message ", string(payload))
				message, err := UnMarshallMessageIn(payload)
				if err != nil {
					slog.Error("-> client : error unMarshalling the payload")
				}
				utils.PrettyDisplay("message in : ", message)
				c.HandleMessage(message)
			}
		}
	}()
}

func (c *Client) GoWrite() {
	go func() {
		defer func() {
			c.manager.RemoveClient(c)
		}()

		c.conn.WriteMessage(websocket.TextMessage, []byte("readyToCommunicate :-)"))

		for {
			message, ok := <-c.egress
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed:", err)
				}
				break
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("failed to send message: %v", err)
			}
			log.Println("message sent")
		}
	}()
}

func (c *Client) HandleMessage(msg MessageIn) {

}
