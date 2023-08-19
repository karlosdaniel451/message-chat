package model

import (
	"gorm.io/gorm"
)

type CreateGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Group struct {
	gorm.Model
	Name                  string         `json:"name" gorm:"varchar(200); not null"`
	Description           string         `json:"description" gorm:"varchar(1000); not null"`
	ReceivedGroupMessages []GroupMessage `json:"received_group_messages" gorm:"foreignKey:group_id"`
	Users                 []*User        `json:"users" gorm:"many2many:users_groups;"`
}
