package usecase

import (
	"errors"
	"fmt"

	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/errs"
	"github.com/karlosdaniel451/message-chat/repository"
)

type UserUseCase interface {
	Create(user *model.User) (*model.User, error)
	GetById(id uint) (*model.User, error)
	GetByName(name string) ([]*model.User, error)
	GetByEmailAddress(emailAddress string) (*model.User, error)
	DeleteById(id uint) error
	GetAll() ([]*model.User, error)

	SendMessageToUser(message *model.PrivateMessage) (*model.PrivateMessage, error)
	SendMessageToGroup(message *model.GroupMessage) (*model.GroupMessage, error)
}

type UserUseCaseImpl struct {
	repository            repository.UserRepository
	privateMessageUseCase PrivateMessageUseCase
	groupMessageUseCase   GroupMessageUseCase
	groupUseCase GroupUseCase
}

func NewUserUseCaseImpl(
	repository repository.UserRepository,
	privateMessageUseCase PrivateMessageUseCase,
	groupMessageUseCase GroupMessageUseCase,
	groupUseCase GroupUseCase,
) UserUseCaseImpl {

	return UserUseCaseImpl{
		repository:            repository,
		privateMessageUseCase: privateMessageUseCase,
		groupMessageUseCase:   groupMessageUseCase,
		groupUseCase:   groupUseCase,
	}
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

func (useCase UserUseCaseImpl) SendMessageToUser(
	message *model.PrivateMessage,
) (*model.PrivateMessage, error) {

	_, err := useCase.GetById(message.ReceiverId)
	if err != nil {
		if errors.As(err, &errs.NotFoundError{}) {
			return nil, errs.NotFoundError{Message: "receiver user not found"}
		}
	}

	return useCase.privateMessageUseCase.Create(message)
}

func (useCase UserUseCaseImpl) SendMessageToGroup(
	message *model.GroupMessage,
) (*model.GroupMessage, error) {

	_, err := useCase.groupUseCase.GetById(message.GroupId)
	if err != nil {
		if errors.As(err, &errs.NotFoundError{}) {
			return nil, errs.NotFoundError{
				Message: fmt.Sprintf("group with id %d not found", message.GroupId),
			}
		}
	}

	return useCase.groupMessageUseCase.Create(message)
}
