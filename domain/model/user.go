package model

import (
	"fmt"

	"gorm.io/gorm"
)

type CreateUser struct {
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
}

type User struct {
	gorm.Model
	Name                    string           `json:"name" gorm:"varchar(200); not null"`
	EmailAddress            string           `json:"email_address" gorm:"varchar(320); not null; uniqueIndex"`
	SentPrivateMessages     []PrivateMessage `json:"sent_private_messages" gorm:"foreignKey:sender_id"`
	ReceivedPrivateMessages []PrivateMessage `json:"received_private_messages" gorm:"foreignKey:receiver_id"`
	SentGroupMessages       []GroupMessage   `json:"sent_group_messages" gorm:"foreignKey:sender_id"`
	Groups                  []*Group         `json:"groups" gorm:"many2many:users_groups;"`
}

func (user User) String() string {
	return fmt.Sprintf(
		"id: %d, name: %s, emailAddress: %s",
		user.ID, user.Name, user.EmailAddress,
	)
}
