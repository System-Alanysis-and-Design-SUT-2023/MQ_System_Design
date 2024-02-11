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

const errorMessage = `expected "%v", got "%v"`

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
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
}

func TestSubscribeAndUnsubscribe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandlerFunc))
	defer server.Close()

	subscriber := models.NewSubscriber()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tests := []struct {
		name       string
		fn         func(*websocket.Conn) error
		err        error
		subscriber bool
	}{
		{
			name:       "subscribe",
			fn:         subscriber.Subscribe,
			err:        nil,
			subscriber: true,
		},
		{
			name:       "subscribe already exists",
			fn:         subscriber.Subscribe,
			err:        models.ErrSubscriberAlreadyExists,
			subscriber: true,
		},
		{
			name:       "unsubscribe",
			fn:         subscriber.Unsubscribe,
			err:        nil,
			subscriber: false,
		},
		{
			name:       "unsubscribe does not exist",
			fn:         subscriber.Unsubscribe,
			err:        models.ErrSubscriberDoesNotExist,
			subscriber: false,
		},
		{
			name:       "subscribe again",
			fn:         subscriber.Subscribe,
			err:        nil,
			subscriber: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(conn); err != tt.err {
				t.Errorf(errorMessage, tt.err, err)
			}

			if tt.subscriber != subscriber.HasSubscriber() {
				if tt.subscriber {
					t.Errorf("expected to have subscriber")
				} else {
					t.Errorf("expected no subscriber")
				}
			}
		})
	}
}

func TestSendData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandlerFunc))
	defer server.Close()

	subscriber := models.NewSubscriber()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if err = subscriber.Subscribe(conn); err != nil {
		t.Errorf(errorMessage, nil, err)
	}
	if !subscriber.HasSubscriber() {
		t.Error("expected to have subscriber")
	}

	t.Run("send data to subscriber", func(t *testing.T) {
		if err := subscriber.SendData(models.NewData("key", "value", 10)); err != nil {
			t.Errorf(errorMessage, nil, err)
		}

		if _, response, err := conn.ReadMessage(); err != nil {
			t.Errorf(errorMessage, nil, err)
		} else if string(response) != `{"Key":"key","Value":"value","Index":10}` {
			t.Errorf("expected %s, got %s", `{"Key":"key","Value":"value","Index":10}`, string(response))
		}
	})

	if err := subscriber.Unsubscribe(conn); err != nil {
		t.Errorf(errorMessage, nil, err)
	}
	if subscriber.HasSubscriber() {
		t.Error("expected no subscriber")
	}

	t.Run("send data to no subscriber", func(t *testing.T) {
		if err := subscriber.SendData(models.NewData("key", "value", 20)); err != models.ErrNoSubscriberExists {
			t.Errorf(errorMessage, models.ErrNoSubscriberExists, err)
		}
	})
}
