package messages

import "context"

type IMessages interface {
	Save(ctx context.Context, message *Message) error
	GetAllMessagesInRoom(ctx context.Context, roomId string) ([]*Message, error)
}
