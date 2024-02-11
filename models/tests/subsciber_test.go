package tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func TestSubscriber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}
	}))
	defer server.Close()

	subscriber := models.NewSubscriber()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tests := []struct {
		name string
		fn   func(*websocket.Conn) error
		err  error
	}{
		{
			name: "subscribe",
			fn:   subscriber.Subscribe,
			err:  nil,
		},
		{
			name: "subscribe already exists",
			fn:   subscriber.Subscribe,
			err:  models.ErrSubscriberAlreadyExists,
		},
		{
			name: "unsubscribe",
			fn:   subscriber.Unsubscribe,
			err:  nil,
		},
		{
			name: "unsubscribe does not exist",
			fn:   subscriber.Unsubscribe,
			err:  models.ErrSubscriberDoesNotExist,
		},
		{
			name: "subscribe again",
			fn:   subscriber.Subscribe,
			err:  nil,
		},
	}

	const errorMessage = "expected %v, got %v"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn(conn)
			if err != tt.err {
				t.Errorf(errorMessage, tt.err, err)
			}
		})
	}

	checkSubscriberStatus(t, subscriber)

	t.Run("send data to subscriber", func(t *testing.T) {
		err := subscriber.SendData(models.NewData("key", "value", 10))
		if err != nil {
			t.Errorf(errorMessage, nil, err)
		}

		_, response, err := conn.ReadMessage()
		if err != nil {
			t.Errorf(errorMessage, nil, err)
		}

		if string(response) != `{"Key":"key","Value":"value","Index":10}` {
			t.Errorf("expected %s, got %s", `{"Key":"key","Value":"value","Index":10}`, string(response))
		}
	})

	subscriber.Unsubscribe(conn)
	checkNoSubscriber(t, subscriber)

	t.Run("send data to no subscriber", func(t *testing.T) {
		err := subscriber.SendData(models.NewData("key", "value", 20))
		if err != models.ErrNoSubscriberExists {
			t.Errorf(errorMessage, models.ErrNoSubscriberExists, err)
		}
	})
}

func checkSubscriberStatus(t *testing.T, subscriber *models.Subscriber) {
	if !subscriber.HasSubscriber() {
		t.Error("expected to have subscriber")
	}
}

func checkNoSubscriber(t *testing.T, subscriber *models.Subscriber) {
	if subscriber.HasSubscriber() {
		t.Error("expected no subscriber")
	}
}
