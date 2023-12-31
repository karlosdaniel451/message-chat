package usecase

import (
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/repository"
)

type GroupUseCase interface {
	Create(group *model.Group) (*model.Group, error)
	GetById(id uint) (*model.Group, error)
	GetByName(name string) (*model.Group, error)
	DeleteById(id uint) error
	GetAll() ([]*model.Group, error)
}

type GroupUseCaseImpl struct {
	repository repository.GroupRepository
}

func NewGroupUseCaseImpl(
	repository repository.GroupRepository,
) GroupUseCaseImpl {

	return GroupUseCaseImpl{repository: repository}
}

func (useCase GroupUseCaseImpl) Create(
	group *model.Group,
) (*model.Group, error) {

	return useCase.repository.Create(group)
}

func (useCase GroupUseCaseImpl) GetById(id uint) (*model.Group, error) {
	return useCase.repository.GetById(id)
}

func (useCase GroupUseCaseImpl) GetByName(name string) (*model.Group, error) {
	return useCase.repository.GetByName(name)
}

func (useCase GroupUseCaseImpl) DeleteById(id uint) error {
	return useCase.repository.DeleteById(id)
}

func (useCase GroupUseCaseImpl) GetAll() ([]*model.Group, error) {
	return useCase.repository.GetAll()
}
