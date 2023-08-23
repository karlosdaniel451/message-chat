package conversions

import (
	"strconv"

	"github.com/karlosdaniel451/message-chat/api/protobuf"
	"github.com/karlosdaniel451/message-chat/domain/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func ModelPrivateMessagetoProto(
	modelMessage *model.PrivateMessage,
) (*protobuf.PrivateMessage, error) {

	protoMessage := protobuf.PrivateMessage{
		Id:          strconv.FormatUint(uint64(modelMessage.ID), 10),
		SenderId:    strconv.FormatUint(uint64(modelMessage.SenderId), 10),
		ReceiverId:  strconv.FormatUint(uint64(modelMessage.ReceiverId), 10),
		TextContent: modelMessage.TextContent,
		CreatedAt:   timestamppb.New(modelMessage.CreatedAt),
		UpdatedAt:   timestamppb.New(modelMessage.UpdatedAt),
		DeletedAt:   timestamppb.New(modelMessage.DeletedAt.Time),
	}
	return &protoMessage, nil
}

func ModelGroupMessagetoProto(
	modelMessage *model.GroupMessage,
) (*protobuf.GroupMessage, error) {

	protoMessage := protobuf.GroupMessage{
		Id:          strconv.FormatUint(uint64(modelMessage.ID), 10),
		SenderId:    strconv.FormatUint(uint64(modelMessage.SenderId), 10),
		GroupId:     strconv.FormatUint(uint64(modelMessage.GroupId), 10),
		TextContent: modelMessage.TextContent,
		CreatedAt:   timestamppb.New(modelMessage.CreatedAt),
		UpdatedAt:   timestamppb.New(modelMessage.UpdatedAt),
		DeletedAt:   timestamppb.New(modelMessage.DeletedAt.Time),
	}
	return &protoMessage, nil
}

func ProtoPrivateMessageToModel(
	protoModel *protobuf.PrivateMessage,
) (*model.PrivateMessage, error) {

	messageId, err := strconv.ParseUint(protoModel.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	senderId, err := strconv.ParseUint(protoModel.SenderId, 10, 64)
	if err != nil {
		return nil, err
	}

	receiverId, err := strconv.ParseUint(protoModel.ReceiverId, 10, 64)
	if err != nil {
		return nil, err
	}

	// Create the Private Message model.
	model := model.PrivateMessage{
		Model: gorm.Model{
			ID:        uint(messageId),
			CreatedAt: protoModel.CreatedAt.AsTime(),
			UpdatedAt: protoModel.UpdatedAt.AsTime(),
			DeletedAt: gorm.DeletedAt{Time: protoModel.DeletedAt.AsTime()},
		},
		TextContent: protoModel.GetTextContent(),
		ReceiverId:  uint(receiverId),
		SenderId:    uint(senderId),
	}

	return &model, nil
}

func ProtoGroupMessageToModel(
	protoModel *protobuf.GroupMessage,
) (*model.GroupMessage, error) {

	messageId, err := strconv.ParseUint(protoModel.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	senderId, err := strconv.ParseUint(protoModel.SenderId, 10, 64)
	if err != nil {
		return nil, err
	}

	groupId, err := strconv.ParseUint(protoModel.GroupId, 10, 64)
	if err != nil {
		return nil, err
	}

	// Create the Private Message model.
	model := model.GroupMessage{
		Model: gorm.Model{
			ID:        uint(messageId),
			CreatedAt: protoModel.CreatedAt.AsTime(),
			UpdatedAt: protoModel.UpdatedAt.AsTime(),
			DeletedAt: gorm.DeletedAt{Time: protoModel.DeletedAt.AsTime()},
		},
		TextContent: protoModel.GetTextContent(),
		SenderId:    uint(senderId),
		GroupId:     uint(groupId),
	}

	return &model, nil
}
