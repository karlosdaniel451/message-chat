package main

import (
	"log"

	"github.com/karlosdaniel451/message-chat/broker"
	"github.com/karlosdaniel451/message-chat/db"
	"github.com/karlosdaniel451/message-chat/repository"
	"github.com/karlosdaniel451/message-chat/usecase"
)

var (
	groupMessageRepository   repository.GroupMessageRepository
	groupRepository          repository.GroupMessageRepository
	privateMessageRepository repository.PrivateMessageRepository
	userRepository           repository.UserRepository

	groupMessageUseCase   usecase.GroupMessageUseCase
	groupUseCase          usecase.UserUseCase
	privateMessageUseCase usecase.PrivateMessageUseCase
	userUseCase           usecase.UserUseCase
)

func setup() {
	assertInterfaces()

	err := db.Connect()
	if err != nil {
		log.Fatalf("error when connecting to database: %s", err)
	}

	err = broker.Connect()
	if err != nil {
		log.Fatalf("error when connecting to NATS server: %s", err)
	}

	groupMessageRepository = repository.NewGroupMessageDB(db.DB)
	groupRepository = repository.NewGroupRepositoryDB(db.DB)
	privateMessageRepository = repository.NewPrivateMessageDB(db.DB)
	userRepository = repository.NewUserRepositoryDB(db.DB)

	groupMessageUseCase = usecase.NewGroupMessageUseCaseImpl(groupMessageRepository)
	groupRepository = usecase.NewGroupMessageUseCaseImpl(groupRepository)
	privateMessageUseCase = usecase.NewPrivateMessageUseCaseImpl(privateMessageRepository)
	userUseCase = usecase.NewUserUseCaseImpl(
		userRepository,
		privateMessageUseCase,
		groupMessageUseCase,
	)
}

func assertInterfaces() {
	var _ usecase.GroupMessageUseCase = usecase.GroupMessageUseCaseImpl{}
	var _ repository.GroupMessageRepository = repository.GroupMessageRepositoryDB{}

	var _ usecase.GroupUseCase = usecase.GroupUseCaseImpl{}
	var _ repository.GroupRepository = repository.GroupRepositoryDB{}

	var _ usecase.PrivateMessageUseCase = usecase.PrivateMessageUseCaselImpl{}
	var _ repository.PrivateMessageRepository = repository.PrivateMessageRepositoryDB{}

	var _ usecase.UserUseCase = usecase.UserUseCaseImpl{}
	var _ repository.UserRepository = repository.UserRepositoryDB{}
}
