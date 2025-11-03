package client

import "context"

type Messages interface {
	Save(ctx context.Context, message Message) error
	FindByRoom(ctx context.Context, roomID string) []Message
}
