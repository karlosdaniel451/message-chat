package pubsubcontroller

import (
	"context"
	"fmt"
	"log"
	"strconv"

	conversions "github.com/karlosdaniel451/message-chat/api"
	"github.com/karlosdaniel451/message-chat/api/protobuf"
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/usecase"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type UserController struct {
	natsConn            *nats.Conn
	useCase             usecase.UserUseCase
	groupMessageUseCase usecase.GroupMessageUseCase
}

func NewUserPubController(
	natsConn *nats.Conn,
	useCase usecase.UserUseCase,
	groupMessageUseCase usecase.GroupMessageUseCase,
) *UserController {

	return &UserController{
		natsConn:            natsConn,
		useCase:             useCase,
		groupMessageUseCase: groupMessageUseCase,
	}
}

func (controller *UserController) SendMessageToUser(
	ctx context.Context,
	privateMessage *model.PrivateMessage,
) (*model.PrivateMessage, error) {

	createdPrivateMessage, err := controller.useCase.SendMessageToUser(privateMessage)
	if err != nil {
		return nil, fmt.Errorf("error when inserting private message to database: %s", err)
	}

	protoMessage, err := conversions.ModelPrivateMessagetoProto(privateMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"error when cnoverting private message model to protobuf model: %s",
			err,
		)
	}

	serializedMessage, err := proto.Marshal(protoMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"error when marshalling private message: %s",
			err,
		)
	}

	// The subject name is determined by the receiving User id.
	subject := "user:" + strconv.FormatUint(uint64(privateMessage.ReceiverId), 10)

	// producer, err := controller.natsConn.CreateProducer(pulsar.ProducerOptions{
	// 	Topic: subject,
	// })

	err = controller.natsConn.Publish(subject, serializedMessage)
	if err != nil {
		return nil, fmt.Errorf("error when sending private message to user: %s", err)
	}

	return createdPrivateMessage, nil
}

func (controller *UserController) SendMessageToGroup(
	ctx context.Context,
	groupMessage *model.GroupMessage,
) (*model.GroupMessage, error) {

	createdGroupMessage, err := controller.useCase.SendMessageToGroup(groupMessage)
	if err != nil {
		return nil, fmt.Errorf("error when inserting group message to database: %s", err)
	}

	protoMessage, err := conversions.ModelGroupMessagetoProto(groupMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"error when converting GroupMessage model to protobuf model: %s",
			err,
		)
	}

	serializedMessage, err := proto.Marshal(protoMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"error when serializing GroupMessage protobuf model to bytes: %s",
			err,
		)
	}

	// The subject name is determined by the receiving Group id.
	subject := "group:" + strconv.FormatUint(uint64(groupMessage.GroupId), 10)

	err = controller.natsConn.Publish(subject, serializedMessage)
	if err != nil {
		return nil, fmt.Errorf("error when sending message to group: %s", err)
	}

	return createdGroupMessage, nil
}

// Return a read-only channel with all the messages sent from the User identified by
// `senderId` to the User identified by `receiverId`.
func (controller *UserController) ConnectToUser(
	ctx context.Context, senderId uint, receiverId uint,
) <-chan model.PrivateMessage {

	privateMessagesChannel := make(chan model.PrivateMessage)
	natsMessagesChannel := make(chan *nats.Msg)

	// The subject is determined by the id of the User receiving a message.
	subject := "user:" + strconv.FormatUint(uint64(receiverId), 10)

	subscriber, err := controller.natsConn.ChanSubscribe(subject, natsMessagesChannel)
	if err != nil {
		log.Fatalf("error when subscribing to a NATS subject: %s", err)
	}

	go func() {
		defer func() {
			subscriber.Unsubscribe()
			close(privateMessagesChannel)
		}()

		for {
			// Check for cancellation signal.
			select {
			case <-ctx.Done():
				// I had forgotten to include this return at the first time ;(
				return
			case natsMsg := <-natsMessagesChannel:
				var protobufMessage protobuf.PrivateMessage

				err := proto.Unmarshal(natsMsg.Data, &protobufMessage)
				if err != nil {
					log.Printf("error when converting unmarshalling protobuf: %s",
						err)
					continue
				}

				// Filter the messages received from other users
				if protobufMessage.SenderId != strconv.FormatUint(uint64(senderId), 10) {
					continue
				}

				modelMessage, err := conversions.ProtoPrivateMessageToModel(&protobufMessage)
				if err != nil {
					log.Printf("error when converting protobuf PrivateMessage to Model: %s",
						err)
					continue
				}

				privateMessagesChannel <- *modelMessage
			}
		}
	}()

	return privateMessagesChannel
}

func (controller *UserController) ConnectToGroup(
	senderId uint, groupId uint,
) <-chan model.GroupMessage {

	groupMessagesChannel := make(chan model.GroupMessage)
	natsMessagesChannel := make(chan *nats.Msg)

	// The subject is determined by the id of the Group receiving a message.
	subject := "group:" + strconv.FormatUint(uint64(groupId), 10)

	subscriber, err := controller.natsConn.ChanSubscribe(subject, natsMessagesChannel)
	if err != nil {
		log.Fatalf("error when subscribing to a NATS subject: %s", err)
	}

	go func() {
		defer func() {
			subscriber.Unsubscribe()
			close(groupMessagesChannel)
			log.Print("Group Messages Subscriber channel closed")
		}()

		for natsMsg := range natsMessagesChannel {
			var protobufMessage protobuf.GroupMessage

			err := proto.Unmarshal(natsMsg.Data, &protobufMessage)
			if err != nil {
				log.Printf("error when converting unmarshalling protobuf: %s",
					err)
				continue
			}

			// Filter the messages received from other users
			if protobufMessage.SenderId != strconv.FormatUint(uint64(senderId), 10) {
				continue
			}

			modelMessage, err := conversions.ProtoGroupMessageToModel(&protobufMessage)
			if err != nil {
				log.Printf("error when converting protobuf PrivateMessage to Model: %s",
					err)
				continue
			}

			groupMessagesChannel <- *modelMessage
		}
	}()

	return groupMessagesChannel
}
