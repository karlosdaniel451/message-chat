package model

import "gorm.io/gorm"

type CreateGroupMessage struct {
	TextContent string `json:"text_content"`
	SenderId    uint   `json:"sender_id"`
	GroupId     uint   `json:"group_id"`
}

type GroupMessage struct {
	gorm.Model
	TextContent string `json:"text_content" gorm:"varchar(10000); not null"`
	SenderId    uint   `json:"sender_id"`
	GroupId     uint   `json:"group_id"`
}
