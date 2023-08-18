package repository

import (
	"fmt"

	"github.com/karlosdaniel451/go-rest-api-template/domain/model"
	"github.com/karlosdaniel451/go-rest-api-template/errs"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(user *model.Task) (*model.Task, error)
	GetById(id uint) (*model.Task, error)
	GetByName(name string) ([]*model.Task, error)
	GetByDescription(description string) ([]*model.Task, error)
	DeleteById(id uint) error
	GetAll() ([]*model.Task, error)
}

type TaskRepositoryDB struct {
	db *gorm.DB
}

func NewTaskRepositoryDB(db *gorm.DB) *TaskRepositoryDB {
	return &TaskRepositoryDB{db: db}
}

func (repository TaskRepositoryDB) Create(task *model.Task) (*model.Task, error) {
	result := repository.db.Create(task)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("it was not possible to insert task: %s", result.Error)
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

func (repository TaskRepositoryDB) GetById(id uint) (*model.Task, error) {
	var task model.Task

	result := repository.db.First(&task, "id = ?", id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("there is no task with id %d", id),
			}
		}
		return nil, result.Error
	}

	return &task, nil
}

func (repository TaskRepositoryDB) GetByName(name string) ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)

	result := repository.db.Where("NAME LIKE %?%", name)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (repository TaskRepositoryDB) GetByDescription(description string) ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)

	result := repository.db.Where("description LIKE %?%", description)
	if result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

func (repository TaskRepositoryDB) DeleteById(id uint) error {
	var user model.Task

	result := repository.db.First(&user, id)
	if result.Error != nil {
		if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
			return errs.NotFoundError{
				Message: fmt.Sprintf("there is no task with id %d", id),
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

func (repository TaskRepositoryDB) GetAll() ([]*model.Task, error) {
	allUsers := make([]*model.Task, 0)

	result := repository.db.Find(&allUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	return allUsers, nil
}
