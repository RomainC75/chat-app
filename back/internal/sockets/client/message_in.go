package client

func CreateMessageIn(mType MessageInType, content map[string]string) MessageIn {
	mi := MessageIn{
		Type:    mType,
		Content: content,
	}
	return mi
}

func CreateBroadcastMessageIn(message string) MessageIn {
	return CreateMessageIn(BROADCAST_MESSAGE, map[string]string{
		"message": message,
	})
}
