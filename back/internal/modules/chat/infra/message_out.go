package chat_app_infra

import chat_client "chat/internal/modules/chat/domain/client"

type MessageOut struct {
	Type    chat_client.MessageOutType `json:"type"`
	Content map[string]string          `json:"content"`
}

func BuildMessageOut(mType chat_client.MessageOutType, content map[string]string) MessageOut {
	mo := MessageOut{
		Type:    mType,
		Content: content,
	}
	return mo
}
