package model

import (
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
