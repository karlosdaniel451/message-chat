package model

import (
	"fmt"

	"gorm.io/gorm"
)

type CreatePrivateMessage struct {
	TextContent string `json:"text_content"`
	SenderId    uint   `json:"sender_id"`
	ReceiverId  uint   `json:"receiver_id"`
}

type PrivateMessage struct {
	gorm.Model
	TextContent string `json:"text_content" gorm:"varchar(10000); not null"`
	SenderId    uint   `json:"sender_id"`
	ReceiverId  uint   `json:"receiver_id"`
}

func (privateMessage PrivateMessage) String() string {
	return fmt.Sprintf(
		"id: %d, text content: %s, sender id: %d, receiver id: %d",
		privateMessage.ID,
		privateMessage.TextContent,
		privateMessage.SenderId,
		privateMessage.ReceiverId,
	)
}
