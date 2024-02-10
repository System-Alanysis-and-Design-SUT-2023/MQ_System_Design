package tools

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func HandlePush(repo *Repository, args ...string) string {
	if len(args) < 2 {
		return "Invalid push command!\nUse: push <topic> <message>"
	}

	topic, message := args[0], args[1]
	err := repo.Push(topic, message)
	if err != nil {
		return fmt.Sprintf("Error pushing message: %s", err)
	}

	return "Message pushed successfully!"
}

func HandlePull(repo *Repository) string {
	data, err := repo.Pull()
	if err != nil {
		return fmt.Sprintf("Error pulling message: %s", err)
	}

	return fmt.Sprintf("Message: %s", data)
}

func HandleSubscribe(repo *Repository, conn *websocket.Conn) string {
	err := repo.Subscribe(conn)
	if err != nil {
		return fmt.Sprintf("Error subscribing: %s", err)
	}

	return "Subscribed successfully!"
}

func HandleUnsubscribe(repo *Repository, conn *websocket.Conn) string {
	err := repo.Unsubscribe(conn)
	if err != nil {
		return fmt.Sprintf("Error unsubscribing: %s", err)
	}

	return "Unsubscribed successfully!"
}
