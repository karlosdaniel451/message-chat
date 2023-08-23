package pubsubcontroller

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/karlosdaniel451/message-chat/api/protobuf"
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/usecase"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserController struct {
	natsConn            nats.Conn
	useCase             usecase.UserUseCase
	groupMessageUseCase usecase.GroupMessageUseCase
}

func NewUserController(
	useCase usecase.UserUseCase,
	groupMessageUseCase usecase.GroupMessageUseCase,
) *UserController {

	return &UserController{useCase: useCase, groupMessageUseCase: groupMessageUseCase}
}

func (controller *UserController) SendMessage(
	groupMessage *model.GroupMessage,
) (*model.GroupMessage, error) {

	createdGroupMessage, err := controller.useCase.SendMessageToGroup(groupMessage)

	if err != nil {
		return nil, fmt.Errorf("error when inserting group message to database: %s", err)
	}

	serializedGroupMessage, err := proto.Marshal(&protobuf.GroupMessage{
		Id:          strconv.FormatUint(uint64(createdGroupMessage.ID), 10),
		SenderId:    strconv.FormatUint(uint64(createdGroupMessage.SenderId), 10),
		GroupId:     strconv.FormatUint(uint64(createdGroupMessage.GroupId), 10),
		TextContent: createdGroupMessage.TextContent,
		CreatedAt:   timestamppb.New(createdGroupMessage.CreatedAt),
		UpdatedAt:   timestamppb.New(createdGroupMessage.UpdatedAt),
		DeletedAt:   timestamppb.New(createdGroupMessage.DeletedAt.Time),
	})

	if err != nil {
		return nil, fmt.Errorf("error when serializing group message to protobuf: %s", err)
	}

	err = controller.natsConn.Publish(
		strconv.FormatUint(uint64(createdGroupMessage.GroupId), 10),
		serializedGroupMessage,
	)

	if err != nil {
		return nil, fmt.Errorf("error when publishing group message: %s", err)
	}

	return createdGroupMessage, nil
}

/*
func (controller *UserController) ConnectToUser(
	senderId uint, receiverId uint,
) chan model.PrivateMessage {

	natsMsgChannel := make(chan *nats.Msg)
	privateMessagesChan := make(chan model.PrivateMessage)

	controller.natsConn.ChanSubscribe(
		strconv.FormatUint(uint64(receiverId), 10), natsMsgChannel,
	)

	go func() {
		for natsMsg := range natsMsgChannel {
			serializedMessage := protobuf.PrivateMessage{}

			proto.Unmarshal(natsMsg.Data, &serializedMessage)

			// Filter the messages received from other users
			if serializedMessage.SenderId != strconv.FormatUint(uint64(senderId), 10){
				continue
			}

			messageId, err := strconv.ParseUint(serializedMessage.Id, 10, 64)
			if err != nil {
				panic(err)
			}

			senderId, err := strconv.ParseUint(serializedMessage.SenderId, 10, 64)
			if err != nil {
				panic(err)
			}

			receiverId, err := strconv.ParseUint(serializedMessage.ReceiverId, 10, 64)
			if err != nil {
				panic(err)
			}

			groupMessage := model.PrivateMessage{
				Model: gorm.Model{
					ID:        uint(messageId),
					CreatedAt: serializedMessage.CreatedAt.AsTime(),
					UpdatedAt: serializedMessage.UpdatedAt.AsTime(),
					DeletedAt: gorm.DeletedAt{Time: serializedMessage.DeletedAt.AsTime()},
				},
				TextContent: serializedMessage.GetTextContent(),
				ReceiverId:  uint(receiverId),
				SenderId:    uint(senderId),
			}
			privateMessagesChan <- groupMessage
		}
	}()

	return privateMessagesChan
}
*/

func (controller *UserController) ConnectToUser(
	senderId uint, receiverId uint,
) chan model.PrivateMessage {

	privateMessagesChan := make(chan model.PrivateMessage)

	sub, err := controller.natsConn.SubscribeSync(
		strconv.FormatUint(uint64(receiverId), 10))

	if err != nil {
		log.Fatalf("error when subscribing to a NATS subject: %s", err)
	}

	go func() {
		for {
			natsMsg, err := sub.NextMsgWithContext(context.Background())
			if err != nil {
				log.Fatalf("error when trying to receive a NATS message: %s", err)
			}

			serializedMessage := protobuf.PrivateMessage{}

			proto.Unmarshal(natsMsg.Data, &serializedMessage)

			// Filter the messages received from other users
			if serializedMessage.SenderId != strconv.FormatUint(uint64(senderId), 10) {
				return
			}

			messageId, err := strconv.ParseUint(serializedMessage.Id, 10, 64)
			if err != nil {
				panic(err)
			}

			senderId, err := strconv.ParseUint(serializedMessage.SenderId, 10, 64)
			if err != nil {
				panic(err)
			}

			receiverId, err := strconv.ParseUint(serializedMessage.ReceiverId, 10, 64)
			if err != nil {
				panic(err)
			}

			// Create the receive Private Message
			receivedMessage := model.PrivateMessage{
				Model: gorm.Model{
					ID:        uint(messageId),
					CreatedAt: serializedMessage.CreatedAt.AsTime(),
					UpdatedAt: serializedMessage.UpdatedAt.AsTime(),
					DeletedAt: gorm.DeletedAt{Time: serializedMessage.DeletedAt.AsTime()},
				},
				TextContent: serializedMessage.GetTextContent(),
				ReceiverId:  uint(receiverId),
				SenderId:    uint(senderId),
			}
			privateMessagesChan <- receivedMessage
		}
	}()

	return privateMessagesChan
}

func (controller *UserController) ConnectToGroup(userId uint) chan model.GroupMessage {
	natsMsgChannel := make(chan *nats.Msg)
	privateMessagesChan := make(chan model.GroupMessage)

	controller.natsConn.ChanSubscribe(
		strconv.FormatUint(uint64(userId), 10), natsMsgChannel,
	)

	go func() {
		for natsMsg := range natsMsgChannel {
			serializedMessage := protobuf.GroupMessage{}

			proto.Unmarshal(natsMsg.Data, &serializedMessage)

			messageId, err := strconv.ParseUint(serializedMessage.GroupId, 10, 64)
			if err != nil {
				panic(err)
			}
			senderId, err := strconv.ParseUint(serializedMessage.SenderId, 10, 64)
			if err != nil {
				panic(err)
			}

			groupMessage := model.GroupMessage{
				Model: gorm.Model{
					ID:        uint(messageId),
					CreatedAt: serializedMessage.CreatedAt.AsTime(),
					UpdatedAt: serializedMessage.UpdatedAt.AsTime(),
					DeletedAt: gorm.DeletedAt{Time: serializedMessage.DeletedAt.AsTime()},
				},
				TextContent: serializedMessage.GetTextContent(),
				GroupId:     userId,
				SenderId:    uint(senderId),
			}
			privateMessagesChan <- groupMessage
		}
	}()

	return privateMessagesChan
}
