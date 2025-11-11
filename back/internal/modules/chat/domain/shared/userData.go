package chat_shared

import "github.com/google/uuid"

type UserData struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
