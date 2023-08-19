package usecase

import (
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/repository"
)

type GroupMessageUseCase interface {
	Create(message *model.GroupMessage) (*model.GroupMessage, error)
	GetById(id uint) (*model.GroupMessage, error)
	DeleteById(id uint) error
	GetAll() ([]*model.GroupMessage, error)
}

type GroupMessageUseCaseImpl struct {
	repository repository.GroupMessageRepository
}

func NewGroupMessageUseCaseImpl(
	repository repository.GroupMessageRepository,
) GroupMessageUseCaseImpl {

	return GroupMessageUseCaseImpl{repository: repository}
}

func (useCase GroupMessageUseCaseImpl) Create(
	message *model.GroupMessage,
) (*model.GroupMessage, error) {

	return useCase.repository.Create(message)
}

func (useCase GroupMessageUseCaseImpl) GetById(id uint) (*model.GroupMessage, error) {
	return useCase.repository.GetById(id)
}

func (useCase GroupMessageUseCaseImpl) DeleteById(id uint) error {
	return useCase.repository.DeleteById(id)
}

func (useCase GroupMessageUseCaseImpl) GetAll() ([]*model.GroupMessage, error) {
	return useCase.repository.GetAll()
}
