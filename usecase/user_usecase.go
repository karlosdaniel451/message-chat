package usecase

import (
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/repository"
)

type UserUseCase interface {
	Create(user *model.User) (*model.User, error)
	GetById(id uint) (*model.User, error)
	GetByName(name string) ([]*model.User, error)
	GetByEmailAddress(emailAddress string) (*model.User, error)
	DeleteById(id uint) error
	GetAll() ([]*model.User, error)

	SendMessageToGroup(message *model.GroupMessage) (*model.GroupMessage, error)
}

type UserUseCaseImpl struct {
	repository repository.UserRepository
}

func NewUserUseCaseImpl(
	repository repository.UserRepository,
) UserUseCaseImpl {

	return UserUseCaseImpl{repository: repository}
}

func (useCase UserUseCaseImpl) Create(user *model.User) (*model.User, error) {
	return useCase.repository.Create(user)
}

func (useCase UserUseCaseImpl) GetById(id uint) (*model.User, error) {
	return useCase.repository.GetById(id)
}

func (useCase UserUseCaseImpl) GetByName(name string) ([]*model.User, error) {
	return useCase.repository.GetByName(name)
}

func (useCase UserUseCaseImpl) GetByEmailAddress(emailAddress string) (*model.User, error) {
	return useCase.repository.GetByEmailAddress(emailAddress)
}

func (useCase UserUseCaseImpl) DeleteById(id uint) error {
	return useCase.repository.DeleteById(id)
}

func (useCase UserUseCaseImpl) GetAll() ([]*model.User, error) {
	return useCase.repository.GetAll()
}

func (useCase UserUseCaseImpl) SendMessageToGroup(
	message *model.GroupMessage,
	groupMessageRepository repository.GroupMessageRepository,
) (*model.GroupMessage, error) {

	return groupMessageRepository.Create(message)
}
