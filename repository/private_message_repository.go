package repository

import (
	"fmt"

	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/errs"
	"gorm.io/gorm"
)

type PrivateMessageRepository interface {
	Create(message *model.PrivateMessage) (*model.PrivateMessage, error)
	GetById(id uint) (*model.PrivateMessage, error)
	DeleteById(id uint) error
	GetChatConversation(senderId, receiverId uint) ([]*model.PrivateMessage, error)
	GetAll() ([]*model.PrivateMessage, error)
}

type PrivateMessageRepositoryDB struct {
	db *gorm.DB
}

func NewPrivateMessageDB(db *gorm.DB) *PrivateMessageRepositoryDB {
	return &PrivateMessageRepositoryDB{db: db}
}

func (repository PrivateMessageRepositoryDB) Create(
	message *model.PrivateMessage,
) (*model.PrivateMessage, error) {

	result := repository.db.Create(message)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf(
			"error when inserting private message: %s",
			result.Error,
		)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

func (repository PrivateMessageRepositoryDB) GetById(id uint) (*model.PrivateMessage, error) {
	var message model.PrivateMessage

	result := repository.db.First(&message, "id = ?", id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no private message with id %d", id),
			}
		}
		return nil, result.Error
	}

	return &message, nil
}

func (repository PrivateMessageRepositoryDB) DeleteById(id uint) error {
	var message model.PrivateMessage

	result := repository.db.First(&message, id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return errs.NotFoundError{
				Message: fmt.Sprintf("there is no private message with id %d", id),
			}
		}
		return result.Error
	}
	result = result.Delete(&message)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository PrivateMessageRepositoryDB) GetChatConversation(
	senderId uint,
	receiverId uint,
) ([]*model.PrivateMessage, error) {

	chatPrivateMessages := make([]*model.PrivateMessage, 0)

	// Filter chat messages.
	result := repository.db.
		Where(
			"sender_id = ? AND receiver_id = ? OR sender_id = ? AND receiver_id = ?",
			senderId, receiverId, receiverId, senderId,
		).
		Find(&chatPrivateMessages)

	if result.Error != nil {
		return nil, result.Error
	}

	return chatPrivateMessages, nil
}

func (repository PrivateMessageRepositoryDB) GetAll() ([]*model.PrivateMessage, error) {
	allPrivateMessages := make([]*model.PrivateMessage, 0)

	result := repository.db.Find(&allPrivateMessages)
	if result.Error != nil {
		return nil, result.Error
	}

	return allPrivateMessages, nil
}
