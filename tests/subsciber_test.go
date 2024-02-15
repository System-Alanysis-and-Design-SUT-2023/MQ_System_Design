package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
	"github.com/gorilla/websocket"
)

func TestSubscribeAndUnsubscribe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandlerFunc))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	subscriber := models.NewSubscriber()
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

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	subscriber := models.NewSubscriber()
	if err = subscriber.Subscribe(conn); err != nil {
		t.Errorf(errorMessage, nil, err)
	}
	if !subscriber.HasSubscriber() {
		t.Error("expected to have subscriber")
	}

	t.Run("send data to subscriber", func(t *testing.T) {
		data := models.NewData("USA", "California", 10)
		if err := subscriber.SendData(data); err != nil {
			t.Errorf(errorMessage, nil, err)
		}

		expectedMessage := fmt.Sprintf(`["%s","%s"]`, data.Key, data.Value)
		if _, response, err := conn.ReadMessage(); err != nil {
			t.Errorf(errorMessage, nil, err)
		} else if string(response) != expectedMessage {
			t.Errorf(errorMessageString, expectedMessage, string(response))
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
