package repository

import (
	"fmt"

	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/errs"
	"gorm.io/gorm"
)

type GroupMessageRepository interface {
	Create(message *model.GroupMessage) (*model.GroupMessage, error)
	GetById(id uint) (*model.GroupMessage, error)
	DeleteById(id uint) error
	GetAll() ([]*model.GroupMessage, error)
}

type GroupMessageRepositoryDB struct {
	db *gorm.DB
}

func NewGroupMessageDB(db *gorm.DB) *GroupMessageRepositoryDB {
	return &GroupMessageRepositoryDB{db: db}
}

func (repository GroupMessageRepositoryDB) Create(
	message *model.GroupMessage,
) (*model.GroupMessage, error) {

	result := repository.db.Create(message)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf(
			"error when inserting group message: %s",
			result.Error,
		)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return message, nil
}

func (repository GroupMessageRepositoryDB) GetById(id uint) (*model.GroupMessage, error) {
	var message model.GroupMessage

	result := repository.db.First(&message, "id = ?", id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no group message with id %d", id),
			}
		}
		return nil, result.Error
	}

	return &message, nil
}

func (repository GroupMessageRepositoryDB) DeleteById(id uint) error {
	var message model.GroupMessage

	result := repository.db.First(&message, id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return errs.NotFoundError{
				Message: fmt.Sprintf("there is no group message with id %d", id),
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

func (repository GroupMessageRepositoryDB) GetAll() ([]*model.GroupMessage, error) {
	allGroupMessages := make([]*model.GroupMessage, 0)

	result := repository.db.Find(&allGroupMessages)
	if result.Error != nil {
		return nil, result.Error
	}

	return allGroupMessages, nil
}
