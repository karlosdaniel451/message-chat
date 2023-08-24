package usecase

import (
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/repository"
)

type PrivateMessageUseCase interface {
	Create(message *model.PrivateMessage) (*model.PrivateMessage, error)
	GetById(id uint) (*model.PrivateMessage, error)
	DeleteById(id uint) error
	GetChatConversation(senderId, receiverId uint) ([]*model.PrivateMessage, error)
	GetAll() ([]*model.PrivateMessage, error)
}

type PrivateMessageUseCaselImpl struct {
	repository repository.PrivateMessageRepository
}

func NewPrivateMessageUseCaseImpl(
	repository repository.PrivateMessageRepository,
) PrivateMessageUseCaselImpl {

	return PrivateMessageUseCaselImpl{repository: repository}
}

func (useCase PrivateMessageUseCaselImpl) Create(
	user *model.PrivateMessage,
) (*model.PrivateMessage, error) {

	return useCase.repository.Create(user)
}

func (useCase PrivateMessageUseCaselImpl) GetById(id uint) (*model.PrivateMessage, error) {
	return useCase.repository.GetById(id)
}

func (useCase PrivateMessageUseCaselImpl) DeleteById(id uint) error {
	return useCase.repository.DeleteById(id)
}

func (usecase PrivateMessageUseCaselImpl) GetChatConversation(
	senderId, receiverId uint,
) ([]*model.PrivateMessage, error) {

	return usecase.repository.GetAll()
}

func (useCase PrivateMessageUseCaselImpl) GetAll() ([]*model.PrivateMessage, error) {
	return useCase.repository.GetAll()
}
