package controller

import (
	"fmt"
	"strconv"

	"github.com/karlosdaniel451/message-chat/api/protobuf"
	"github.com/karlosdaniel451/message-chat/domain/model"
	"github.com/karlosdaniel451/message-chat/usecase"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GroupMessageController struct {
	natsConn nats.Conn
	useCase  usecase.GroupMessageUseCase
}

func NewGroupMessageController(useCase usecase.GroupMessageUseCase) *GroupMessageController {
	return &GroupMessageController{useCase: useCase}
}

func (controller *GroupMessageController) SendMessage(
	groupMessage *model.GroupMessage,
) (*model.GroupMessage, error) {

	createdGroupMessage, err := controller.useCase.Create(groupMessage)
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

// func (controller *GroupMessageController) ReceiveMessages(msg *nats.Msg) *model.GroupMessage {
// 	// TODO
// 	controller.natsConn.ChanSubscribe()

// 	return nil
// }
