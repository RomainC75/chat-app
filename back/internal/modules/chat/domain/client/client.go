package chat_client

import (
	"chat/internal/modules/chat/domain/messages"
	socket_shared "chat/internal/modules/chat/domain/shared"
	"context"

	"github.com/google/uuid"
)

type IManager interface {
	RemoveClient(*Client)
	BroadcastMessage(message *messages.Message)
	SendRoomMessage(message *messages.Message)
	CreateRoom(c *Client, roomName string, description string)
	ConnectUserAndRoom(c *Client, roomId uuid.UUID) error
}

type IRoom interface {
	GetId() uuid.UUID
	GetName() string
	GetClients() []socket_shared.UserData
}

type Client struct {
	manager IManager
	room    IRoom
	conn    IWebSocket
	user    socket_shared.UserData
	// egress  chan (*messages.Message)
	// eventEgress chan (IEvents)
	cancelFn context.CancelFunc
	ctx      context.Context
}

func NewClient(manager IManager, conn IWebSocket, userData socket_shared.UserData) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Client{
		manager: manager,
		conn:    conn,
		user:    userData,
		// egress:   make(chan *messages.Message),
		cancelFn: cancel,
		ctx:      ctx,
	}
	// c.GoListen()
	// c.GoWrite()
	c.WriteHelloMessage()
	return c
}

func (c *Client) PrepareToBeDeleted() {
	c.cancelFn()
}

func (c *Client) SendMessageToClient(message *messages.Message) {
	c.conn.WriteTextMessage(message)
}

func (c *Client) SendEventToClient(event IEvents) {
	c.conn.WriteEvent(event)
}

func (c *Client) ListenToMessageIn(commandMessageIn ICommandMessageIn) {
	commandMessageIn.Execute(c)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-c.ctx.Done():
	// 			return
	// 		case ICommandMessageIn := <-c.conn.GetChan():
	// 			ICommandMessageIn.Execute(c)
	// 		}
	// 	}
	// }()
}

// func (c *Client) GoWrite() {
// 	c.writeHelloMessage()
// 	go func() {
// 		defer func() {
// 			c.manager.RemoveClient(c)
// 		}()

// 		for {
// 			select {
// 			case <-c.ctx.Done():
// 				return
// 			case message, ok := <-c.egress:
// 				if !ok {
// 					if err := c.conn.WriteCloseMessage(); err != nil {
// 						log.Println("connection closed:", err)
// 					}
// 					continue
// 				}

// 				if err := c.conn.WriteTextMessage(message); err != nil {
// 					log.Printf("failed to send message: %v\n", err)
// 				}
// 				log.Println("message sent")
// 			}
// 		}
// 	}()
// }

func (c *Client) GetUserData() socket_shared.UserData {
	return c.user
}

// === IN ===

func (c *Client) BroadcastMessage(message *messages.Message) {
	c.manager.BroadcastMessage(message)
}

func (c *Client) SendRoomMessage(message *messages.Message) {
	c.manager.SendRoomMessage(message)
}

func (c *Client) CreateRoom(roomName string, description string) {
	c.manager.CreateRoom(c, roomName, description)
}

func (c *Client) ConnectUserToRoom(roomId uuid.UUID) {
	_ = c.manager.ConnectUserAndRoom(c, roomId)
}

// === OUT ===
func (c *Client) WriteHelloMessage() {
	event := &HelloEvent{}
	c.conn.WriteEvent(event)
}

func (c *Client) ConnectToRoom(room IRoom) {
	event := &ConnectedToRoomEvent{
		Users:    c.room.GetClients(),
		RoomName: c.room.GetName(),
		RoomId:   room.GetId(),
	}
	c.conn.WriteEvent(event)
}

func (c *Client) SendRoomCreatedMessage(room IRoom) {
	event := &NewRoomCreatedEvent{
		roomName: c.room.GetName(),
		roomId:   room.GetId(),
	}
	c.conn.WriteEvent(event)
}
