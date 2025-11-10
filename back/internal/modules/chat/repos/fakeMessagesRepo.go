package chat_repos

import (
	"chat/internal/modules/chat/domain/messages"
	"context"
)

type InMemoryMessagesRepo struct {
	savedMessages []*messages.Message
}

func NewInMemoryMessagesRepo() *InMemoryMessagesRepo {
	return &InMemoryMessagesRepo{
		savedMessages: []*messages.Message{},
	}
}

func (repo *InMemoryMessagesRepo) Save(ctx context.Context, message *messages.Message) error {
	repo.savedMessages = append(repo.savedMessages, message)
	return nil
}

func (repo *InMemoryMessagesRepo) GetAllMessagesInRoom(ctx context.Context, roomId string) ([]*messages.Message, error) {
	messagesInRoom := make([]*messages.Message, 0, 100)
	for _, msg := range repo.savedMessages {
		if msg.RoomID().String() == roomId {
			messagesInRoom = append(messagesInRoom, msg)
		}
	}
	return messagesInRoom, nil
}

func (repo *InMemoryMessagesRepo) GetSavedMessages() []*messages.Message {
	return repo.savedMessages
}
