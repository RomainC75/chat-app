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

func (repo *InMemoryMessagesRepo) GetSavedMessages() []*messages.Message {
	return repo.savedMessages
}
