package repository

import (
	"fmt"

	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/errs"
	"gorm.io/gorm"
)

type GroupRepository interface {
	Create(group *model.Group) (*model.Group, error)
	GetById(id uint) (*model.Group, error)
	GetByName(name string) (*model.Group, error)
	DeleteById(id uint) error
	GetAll() ([]*model.Group, error)
}

type GroupRepositoryDB struct {
	db *gorm.DB
}

func NewGroupRepositoryDB(db *gorm.DB) *GroupRepositoryDB {
	return &GroupRepositoryDB{db: db}
}

func (repository GroupRepositoryDB) Create(
	group *model.Group,
) (*model.Group, error) {

	result := repository.db.Create(group)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf(
			"error when inserting group: %s",
			result.Error,
		)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return group, nil
}

func (repository GroupRepositoryDB) GetById(id uint) (*model.Group, error) {
	var group model.Group

	result := repository.db.First(&group, "id = ?", id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no group with id %d", id),
			}
		}
		return nil, result.Error
	}

	return &group, nil
}

func (repository GroupRepositoryDB) GetByName(name string) (*model.Group, error) {
	var group model.Group

	result := repository.db.First(&group, "name = ?", name)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no group with name %s", name),
			}
		}
		return nil, result.Error
	}

	return &group, nil
}

func (repository GroupRepositoryDB) DeleteById(id uint) error {
	var group model.Group

	result := repository.db.First(&group, id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return errs.NotFoundError{
				Message: fmt.Sprintf("there is no group with id %d", id),
			}
		}
		return result.Error
	}
	result = result.Delete(&group)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository GroupRepositoryDB) GetAll() ([]*model.Group, error) {
	allGroups := make([]*model.Group, 0)

	result := repository.db.Find(&allGroups)
	if result.Error != nil {
		return nil, result.Error
	}

	return allGroups, nil
}
