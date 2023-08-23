package setup

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/karlosdaniel451/message-chat/api/pub-sub/pubsubcontroller"
	"github.com/karlosdaniel451/message-chat/broker"
	"github.com/karlosdaniel451/message-chat/db"
	"github.com/karlosdaniel451/message-chat/repository"
	"github.com/karlosdaniel451/message-chat/usecase"
)

var (
	// Pub-Sub controllers
	UserPubSubController           pubsubcontroller.UserController
	GroupPubSubController          pubsubcontroller.GroupMessageController
	PrivateMessagePubSubController pubsubcontroller.PrivateMessageController
	// GroupMessagPubSubController controller.

	// Repositories
	GroupMessageRepository   repository.GroupMessageRepository
	GroupRepository          repository.GroupRepository
	PrivateMessageRepository repository.PrivateMessageRepository
	UserRepository           repository.UserRepository

	// Use cases
	GroupMessageUseCase   usecase.GroupMessageUseCase
	GroupUseCase          usecase.GroupUseCase
	PrivateMessageUseCase usecase.PrivateMessageUseCase
	UserUseCase           usecase.UserUseCase
)

func Setup() {
	assertInterfaces()

	// Load an .env file and set the key-value pairs as environment variables.
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// Try to connect to the database server.
	err := db.Connect()
	if err != nil {
		log.Fatalf("error when connecting to database: %s", err)
	}

	// Try to connect to the Pub-Sub broker server.
	err = broker.Connect()
	if err != nil {
		log.Fatalf("error when connecting to Apache Pulsar broker server: %s", err)
	}


	GroupMessageRepository = repository.NewGroupMessageDB(db.DB)
	GroupRepository = repository.NewGroupRepositoryDB(db.DB)
	PrivateMessageRepository = repository.NewPrivateMessageDB(db.DB)
	UserRepository = repository.NewUserRepositoryDB(db.DB)

	GroupMessageUseCase = usecase.NewGroupMessageUseCaseImpl(GroupMessageRepository)
	GroupUseCase = usecase.NewGroupUseCaseImpl(GroupRepository)
	PrivateMessageUseCase = usecase.NewPrivateMessageUseCaseImpl(PrivateMessageRepository)
	UserUseCase = usecase.NewUserUseCaseImpl(
		UserRepository,
		PrivateMessageUseCase,
		GroupMessageUseCase,
	)

	UserPubSubController = *pubsubcontroller.NewUserPubController(
		*broker.GetClient(),
		UserUseCase,
		GroupMessageUseCase,
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
