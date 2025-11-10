package chat_client

import (
	"chat/internal/modules/chat/domain/messages"
	chat_shared "chat/internal/modules/chat/domain/shared"
	shared_domain "chat/internal/modules/shared/domain"
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
	GetClients() []chat_shared.UserData
}

type Client struct {
	messages messages.IMessages
	manager  IManager
	room     IRoom
	conn     IWebSocket
	user     chat_shared.UserData
	cancelFn context.CancelFunc
	ctx      context.Context
	uuidGen  shared_domain.UuidGenerator
	clock    shared_domain.Clock
}

func NewClient(
	messages messages.IMessages,
	manager IManager,
	conn IWebSocket,
	userData chat_shared.UserData,
	uuidGen shared_domain.UuidGenerator,
	clock shared_domain.Clock,
) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		messages: messages,
		manager:  manager,
		conn:     conn,
		user:     userData,
		cancelFn: cancel,
		ctx:      ctx,
		uuidGen:  uuidGen,
		clock:    clock,
	}
	conn.LinkToClient(c)
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
}

func (c *Client) GetUserData() chat_shared.UserData {
	return c.user
}

// === IN ===

func (c *Client) BroadcastMessage(message *messages.Message) {
	_ = c.messages.Save(context.Background(), message)
	c.manager.BroadcastMessage(message)
}

func (c *Client) SendRoomMessage(message *messages.Message) {
	_ = c.messages.Save(context.Background(), message)
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
	c.room = room
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

func (c *Client) RemoveClient() {
	c.manager.RemoveClient(c)
}
