package clicontroller

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/karlosdaniel451/message-chat/cmd/setup"
	"github.com/karlosdaniel451/message-chat/domain/model"
)

func LoginAsUser(reader *bufio.Reader) (*model.User, error) {
	// read email address of user
	fmt.Print("email of user: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	email = strings.TrimSuffix(email, "\n")

	return setup.UserUseCase.GetByEmailAddress(email)
}

func CreateUser(reader *bufio.Reader) (newUser *model.User, err error) {
	user := &model.User{}

	// Read name of user
	fmt.Print("name of user: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	name = strings.TrimSuffix(name, "\n")
	user.Name = name

	// Read email address of user
	fmt.Print("email of user: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	email = strings.TrimSuffix(email, "\n")
	user.EmailAddress = email

	return setup.UserUseCase.Create(user)
}

func DeleteUser(reader *bufio.Reader) error {
	fmt.Print("id of user to be deleted: ")
	userIdString, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	userIdString = strings.TrimSuffix(userIdString, "\n")
	userId, err := strconv.ParseUint(userIdString, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}

	return setup.UserRepository.DeleteById(uint(userId))
}

func ListUsers() ([]*model.User, error) {
	return setup.UserUseCase.GetAll()
}

func CreateGroup(reader *bufio.Reader) (newGroup *model.Group, err error) {
	group := &model.Group{}

	// Read name of group
	fmt.Print("name of group: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	name = strings.TrimSuffix(name, "\n")
	group.Name = name

	// Read description of user
	fmt.Print("description of group: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	description = strings.TrimSuffix(description, "\n")
	group.Description = description

	return setup.GroupUseCase.Create(group)
}

func DeleteGroup(reader *bufio.Reader) error {
	fmt.Print("id of group to be deleted: ")
	groupIdString, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	groupIdString = strings.TrimSuffix(groupIdString, "\n")
	groupId, err := strconv.ParseUint(groupIdString, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid group id")
	}

	return setup.GroupUseCase.DeleteById(uint(groupId))
}

func ListGroups() ([]*model.Group, error) {
	return setup.GroupUseCase.GetAll()
}

func SendToUser(
	ctx context.Context,
	reader *bufio.Reader,
	currentUser *model.User,
) (sentMessage *model.PrivateMessage, err error) {

	var receiverId int

	fmt.Print("receiver Id: ")
	_, err = fmt.Scanf("%d", &receiverId)
	if err != nil {
		return nil, err
	}

	fmt.Print("message: ")
	messageContent, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	messageContent = strings.TrimSuffix(messageContent, "\n")

	sentMessage, err = setup.UserPubSubController.SendMessageToUser(
		ctx,
		&model.PrivateMessage{
			TextContent: messageContent,
			SenderId:    currentUser.ID,
			ReceiverId:  uint(receiverId),
		},
	)

	return sentMessage, err
}

func SendToGroup(
	ctx context.Context,
	reader *bufio.Reader,
	currentUser *model.User,
) (sentMessage *model.GroupMessage, err error) {

	var groupId int

	fmt.Print("group id: ")
	_, err = fmt.Scanf("%d", &groupId)
	if err != nil {
		return nil, err
	}

	fmt.Print("message: ")
	messageContent, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	messageContent = strings.TrimSuffix(messageContent, "\n")

	sentMessage, err = setup.UserPubSubController.SendMessageToGroup(
		ctx,
		&model.GroupMessage{
			TextContent: messageContent,
			SenderId:    currentUser.ID,
			GroupId:     uint(groupId),
		},
	)

	// sentMessage, err = setup.UserUseCase.SendMessageToGroup(&model.GroupMessage{
	// 	TextContent: messageContent,
	// 	SenderId:    uint(currentUser.ID),
	// 	GroupId:     uint(groupId),
	// })

	return sentMessage, err
}

func ConnectToUser(
	ctx context.Context, reader *bufio.Reader, currentUser *model.User,
) (<-chan model.PrivateMessage, *model.User, error) {

	// Read email address of user to be connected to.
	fmt.Print("email of user to be connected to: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, nil, err
	}
	email = strings.TrimSuffix(email, "\n")

	userToBeConnectedTo, err := setup.UserUseCase.GetByEmailAddress(email)
	if err != nil {
		return nil, nil, err
	}

	// Retrieve past messages.
	pastMessages, err := setup.PrivateMessageUseCase.GetChatMessages(
		currentUser.ID, userToBeConnectedTo.ID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error when retrieving past messages: %s", err)
	}

	// Print past messages.
	for _, message := range pastMessages {
		fmt.Printf("%d: %s\n", message.SenderId, message.TextContent)
	}

	privateMessagesChan := setup.UserPubSubController.ConnectToUser(
		ctx, userToBeConnectedTo.ID, currentUser.ID,
	)

	return privateMessagesChan, userToBeConnectedTo, nil
}

func ConnectToGroup(
	reader *bufio.Reader,
	currentUser *model.User,
) (<-chan model.GroupMessage, error) {

	// read name of group to be connected to
	fmt.Print("name of group to be connected to: ")
	groupName, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	groupName = strings.TrimSuffix(groupName, "\n")

	groupToBeConnectedTo, err := setup.GroupUseCase.GetByName(groupName)
	if err != nil {
		return nil, fmt.Errorf("error when retrieving group: %s", err)
	}

	// Print past messages.
	pastMessages := groupToBeConnectedTo.ReceivedGroupMessages
	for _, message := range pastMessages {
		fmt.Printf("%d: %s\n", message.SenderId, message.TextContent)
	}

	groupMessagesChan := setup.UserPubSubController.ConnectToGroup(
		currentUser.ID, groupToBeConnectedTo.ID,
	)

	return groupMessagesChan, nil
}
