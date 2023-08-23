package model

import (
	"fmt"

	"gorm.io/gorm"
)

type CreateGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Group struct {
	gorm.Model
	Name                  string         `json:"name" gorm:"varchar(200); not null, uniqueIndex"`
	Description           string         `json:"description" gorm:"varchar(1000); not null"`
	ReceivedGroupMessages []GroupMessage `json:"received_group_messages" gorm:"foreignKey:group_id"`
	Users                 []*User        `json:"users" gorm:"many2many:users_groups;"`
}

func (group Group) String() string {
	return fmt.Sprintf(
		"id: %d, name: %s, description: %s",
		group.ID, group.Name, group.Description,
	)
}
