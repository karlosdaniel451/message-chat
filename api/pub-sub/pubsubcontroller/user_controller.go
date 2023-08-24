package pubsubcontroller

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/apache/pulsar-client-go/pulsar"
	conversions "github.com/karlosdaniel451/message-chat/api"
	"github.com/karlosdaniel451/message-chat/api/protobuf"
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/usecase"
	"google.golang.org/protobuf/proto"
)

type UserController struct {
	pulsarClient        pulsar.Client
	useCase             usecase.UserUseCase
	groupMessageUseCase usecase.GroupMessageUseCase
}

func NewUserPubController(
	pulsarClient pulsar.Client,
	useCase usecase.UserUseCase,
	groupMessageUseCase usecase.GroupMessageUseCase,
) *UserController {

	return &UserController{
		pulsarClient:        pulsarClient,
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

	marshalled, err := proto.Marshal(protoMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"error when marshalling private message: %s",
			err,
		)
	}

	// The topic name is determined by the receiving User id.
	topic := strconv.FormatUint(uint64(privateMessage.ReceiverId), 10)

	producer, err := controller.pulsarClient.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return nil, fmt.Errorf("error when creating Pulsar producer: %s", err)
	}
	log.Print("producer created sucessfully")

	defer producer.Close()

	_, err = producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: marshalled,
	})
	if err != nil {
		return nil, fmt.Errorf("error when sending private message to user")
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
			"error when serializing convertion GroupMessage model to protobuf model: %s",
			err,
		)
	}

	serializedMessage, err := proto.Marshal(protoMessage)
	if err != nil {
		return nil, fmt.Errorf(
			"error when serializing GroupMessage protobuf model: %s",
			err,
		)
	}

	// The topic name is determined by the receiving Group id.
	topic := strconv.FormatUint(uint64(groupMessage.GroupId), 10)

	producer, err := controller.pulsarClient.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return nil, fmt.Errorf("error when creating Pulsar producer: %s", err)
	}

	defer producer.Close()

	producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: serializedMessage,
	})

	return createdGroupMessage, nil
}

// Return a read-only channel with all the messages sent from the User identified by
// `senderId` to the User identified by `receiverId`.
func (controller *UserController) ConnectToUser(
	senderId uint, receiverId uint,
) <-chan model.PrivateMessage {

	privateMessagesChannel := make(chan model.PrivateMessage)
	pulsarMessagesChannel := make(chan pulsar.ConsumerMessage)

	// The topic is determined by the id of the User receiving a message.
	topic := strconv.FormatUint(uint64(receiverId), 10)

	consumer, err := controller.pulsarClient.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: fmt.Sprintf("%d-%d", senderId, receiverId),
		MessageChannel:   pulsarMessagesChannel,
	})
	if err != nil {
		log.Fatalf("error when subscribing to a Pulsar topic: %s", err)
	}

	// defer consumer.Close()

	go func() {
		// defer close(privateMessagesChannel)
		defer func() {
			consumer.Close()
			close(privateMessagesChannel)
			log.Print("Private Messages Subscriber Channel closed")
		}()

		for pulsarMessage := range pulsarMessagesChannel {
			var protobufMessage protobuf.PrivateMessage

			err := proto.Unmarshal(pulsarMessage.Payload(), &protobufMessage)
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
			consumer.Ack(pulsarMessage.Message)
		}
	}()

	return privateMessagesChannel
}

func (controller *UserController) ConnectToGroup(
	senderId uint, groupId uint,
) <-chan model.GroupMessage {

	groupMessagesChannel := make(chan model.GroupMessage)
	pulsarMessagesChannel := make(chan pulsar.ConsumerMessage)

	// The topic is determined by the id of the Group receiving a message.
	topic := strconv.FormatUint(uint64(groupId), 10)

	consumer, err := controller.pulsarClient.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: fmt.Sprintf("%d-%d", senderId, groupId),
		MessageChannel:   pulsarMessagesChannel,
	})
	if err != nil {
		log.Fatalf("error when subscribing to a Pulsar topic: %s", err)
	}

	defer consumer.Close()

	go func() {
		defer close(groupMessagesChannel)

		for pulsarMessage := range pulsarMessagesChannel {
			var protobufMessage protobuf.GroupMessage

			proto.Unmarshal(pulsarMessage.Payload(), &protobufMessage)

			// // Filter the messages received from other users
			// if protobufMessage.SenderId != strconv.FormatUint(uint64(senderId), 10) {
			// 	return
			// }

			modelMessage, err := conversions.ProtoGroupMessageToModel(&protobufMessage)
			if err != nil {
				log.Printf("error when converting protobuf GroupMessage to Model: %s",
					err)
				continue
			}

			groupMessagesChannel <- *modelMessage
		}
	}()

	return groupMessagesChannel
}
