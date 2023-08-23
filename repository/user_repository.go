package repository

import (
	"fmt"

	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/errs"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	GetById(id uint) (*model.User, error)
	GetByName(fullName string) ([]*model.User, error)
	GetByEmailAddress(emailAddress string) (*model.User, error)
	DeleteById(id uint) error
	GetAll() ([]*model.User, error)
}

type UserRepositoryDB struct {
	db *gorm.DB
}

func NewUserRepositoryDB(db *gorm.DB) *UserRepositoryDB {
	return &UserRepositoryDB{db: db}
}

func (repository UserRepositoryDB) Create(user *model.User) (*model.User, error) {
	result := repository.db.Create(user)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("error when inserting user: %s", result.Error)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repository UserRepositoryDB) GetById(id uint) (*model.User, error) {
	var user model.User

	result := repository.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no user with id %d", id),
			}
		}
		return nil, result.Error
	}

	return &user, nil
}

func (repository UserRepositoryDB) GetByName(name string) ([]*model.User, error) {
	user := make([]*model.User, 0)

	result := repository.db.Where("name LIKE %?%", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repository UserRepositoryDB) GetByEmailAddress(emailAddress string) (*model.User, error) {
	var user model.User

	result := repository.db.First(&user, "email_address = ?", emailAddress)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no user with email %s", emailAddress),
			}
		}
		return nil, result.Error
	}

	return &user, nil
}

func (repository UserRepositoryDB) DeleteById(id uint) error {
	var user model.User

	result := repository.db.First(&user, id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return errs.NotFoundError{
				Message: fmt.Sprintf("there is no user with id %d", id),
			}
		}
		return result.Error
	}
	result = result.Delete(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repository UserRepositoryDB) GetAll() ([]*model.User, error) {
	allUsers := make([]*model.User, 0)

	result := repository.db.Find(&allUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return allUsers, nil
}
