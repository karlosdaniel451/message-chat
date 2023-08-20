package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/karlosdaniel451/message-chat/domain/model"
)

// var action = flag.String("action", "a", "help")
// var userId = flag.Uint("userId", 1, "1")
// var action = flag.String("action", "a", "help")

func main() {
	setup()

	reader := bufio.NewReader(os.Stdin)

	condition := true
	for condition {
		fmt.Print(`action: ("sendUser", "sendGroup", "userMessages", "groupMessages", "exit")`)
		action, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		switch action {
		case "sendUser":
			sentMessage, err := sendToGroup(reader)
			if err != nil {
				log.Printf("error when sending message to user: %s", err)
			}
			log.Printf("message sent: %+v", sentMessage)

		case "sendGroup":
			sentMessage, err := sendToGroup(reader)
			if err != nil {
				log.Printf("error when sending message to group: %s", err)
			}
			log.Printf("message sent: %+v", sentMessage)

		case "userMessages":

		case "groupMessages":

		case "exit":
			condition = false

		default:
			log.Printf("error: invalid option")
			continue
		}
	}
}
func sendToUser(reader *bufio.Reader) (sentMessage *model.PrivateMessage, err error) {
	var senderId int
	var receiverId int

	fmt.Print("senderId: ")
	_, err = fmt.Scanf("%d", &senderId)
	if err != nil {
		return nil, err
	}

	fmt.Print("receiverId: ")
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

	sentMessage, err = userUseCase.SendMessageToUser(&model.PrivateMessage{
		TextContent: messageContent,
		SenderId:    uint(senderId),
		ReceiverId:  uint(receiverId),
	})

	return sentMessage, err
}

func sendToGroup(reader *bufio.Reader) (sentMessage *model.GroupMessage, err error) {
	var senderId int
	var groupId int

	fmt.Print("senderId: ")
	_, err = fmt.Scanf("%d", &senderId)
	if err != nil {
		return nil, err
	}

	fmt.Print("groupId: ")
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

	sentMessage, err = userUseCase.SendMessageToGroup(&model.GroupMessage{
		TextContent: messageContent,
		SenderId:    uint(senderId),
		GroupId:     uint(groupId),
	})

	return sentMessage, err
}
