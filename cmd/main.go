package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// "sync"

	"github.com/karlosdaniel451/message-chat/client/cli/clicontroller"
	"github.com/karlosdaniel451/message-chat/cmd/setup"
	"github.com/karlosdaniel451/message-chat/domain/model"
)

func main() {
	setup.Setup()

	reader := bufio.NewReader(os.Stdin)

	var currentUser *model.User

	condition := true
	for condition {
		fmt.Println("\npossible actions:" +
			"\n - login" +
			"\n - currentUser" +
			"\n - createUser" +
			"\n - listUsers" +
			"\n - deleteUser" +
			"\n - createGroup" +
			"\n - listGroups" +
			"\n - deleteGroup" +
			"\n - sendToUser" +
			"\n - sendToGroup" +
			"\n - connectToUser" +
			"\n - connectToGroup" +
			"\n - exit",
		)
		action, err := reader.ReadString('\n')
		action = strings.TrimSuffix(action, "\n")
		if err != nil {
			log.Fatal(err)
		}

		switch action {
		case "login":
			currentUser, err = clicontroller.LoginAsUser(reader)
			if err != nil {
				fmt.Printf("error when doing login: %s\n", err)
				continue
			}

			fmt.Printf("login done as user: %s\n", currentUser)

		case "currentUser":
			if currentUser == nil {
				fmt.Println("login not done yet")
				continue
			}

			fmt.Printf("login done to user: %s\n", currentUser)

		case "createUser":
			createdUser, err := clicontroller.CreateUser(reader)
			if err != nil {
				log.Printf("error when creating user: %s", err)
				continue
			}
			fmt.Printf("user created: %s\n", createdUser)

		case "listUsers":
			allUsers, err := clicontroller.ListUsers()
			if err != nil {
				log.Printf("error when listing users: %s", err)
				continue
			}

			if len(allUsers) == 0 {
				fmt.Println("no users created yet")
				continue
			}

			fmt.Printf("found %d user(s)\n", len(allUsers))
			for _, user := range allUsers {
				fmt.Println(user)
			}

		case "deleteUser":
			err := clicontroller.DeleteUser(reader)
			if err != nil {
				log.Printf("error when deleting user: %s", err)
				continue
			}

			fmt.Println("user deleted successfully")

		case "createGroup":
			createdGroup, err := clicontroller.CreateGroup(reader)
			if err != nil {
				log.Printf("error when creating group: %s", err)
				continue
			}

			fmt.Printf("group created: %s\n", createdGroup)

		case "listGroups":
			allGroups, err := clicontroller.ListGroups()
			if err != nil {
				log.Printf("error when listing groups: %s", err)
				continue
			}

			if len(allGroups) == 0 {
				fmt.Println("no group created yet")
				continue
			}

			fmt.Printf("found %d groups(s)\n", len(allGroups))
			for _, group := range allGroups {
				fmt.Println(group)
			}

		case "deleteGroup":
			err := clicontroller.DeleteGroup(reader)
			if err != nil {
				log.Printf("error when deleting group: %s", err)
				continue
			}

			fmt.Println("group deleted successfully")

		case "sendToUser":
			if currentUser == nil {
				fmt.Println("login not done yet")
				continue
			}

			sentMessage, err := clicontroller.SendToUser(
				context.Background(),
				reader,
				currentUser,
			)
			if err != nil {
				log.Printf("error when sending message to user: %s", err)
				continue
			}

			log.Printf("message sent: %s", sentMessage)

		case "sendToGroup":
			if currentUser == nil {
				fmt.Println("login not done yet")
				continue
			}

			sentMessage, err := clicontroller.SendToGroup(
				context.Background(),
				reader,
				currentUser,
			)
			if err != nil {
				log.Printf("error when sending message to group: %s", err)
				continue
			}

			log.Printf("message sent: %+v", sentMessage)

		case "connectToUser":
			if currentUser == nil {
				fmt.Println("login not done yet")
				continue
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			receivedMessagesChannel, userConnectedTo, err := clicontroller.ConnectToUser(
				ctx, reader, currentUser,
			)
			if err != nil {
				fmt.Printf("error when connecting to user: %s", err)
				continue
			}

			// var wg sync.WaitGroup
			done := make(chan struct{})

			// Read and send new messages to the connected User.
			// wg.Add(1)
			go func() {
				for {
					newMessageContent, err := reader.ReadString('\n')
					if err != nil {
						fmt.Printf("error when reading new message: %s", err)
						return
					}
					newMessageContent = strings.TrimSuffix(newMessageContent, "\n")

					if newMessageContent == `\exit` {
						break
					}

					_, err = setup.UserPubSubController.SendMessageToUser(
						context.Background(),
						&model.PrivateMessage{
							TextContent: newMessageContent,
							SenderId:    currentUser.ID,
							ReceiverId:  userConnectedTo.ID,
						},
					)
					if err != nil {
						log.Printf("error when sending message to user: %s", err)
						continue
					}
				}
				cancel()
				done <- struct{}{}
			}()

			// Print messages as they are received (consumed).
			for message := range receivedMessagesChannel {
				fmt.Printf("%d: %s\n", message.SenderId, message.TextContent)
			}

			<-done
			fmt.Println("got here!!!")

		case "connectToGroup":
			if currentUser == nil {
				fmt.Println("login not done yet")
				continue
			}

			messagesChan, err := clicontroller.ConnectToGroup(reader, currentUser)
			if err != nil {
				fmt.Printf("error when connecting to group: %s", err)
				continue
			}

			// Print messages as they are received (consumed).
			for message := range messagesChan {
				fmt.Printf("%d: %s\n", message.SenderId, message.TextContent)
			}

		case "exit":
			condition = false

		default:
			fmt.Printf("error: invalid option")
			continue
		}
	}
}
