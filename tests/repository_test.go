package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/internals"
	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/models"
	"github.com/gorilla/websocket"
)

func TestRepositoryActions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(HandlerFunc))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	repository := internals.NewRepository()
	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "push",
			fn: func() {
				err := repository.Push("Iran", "Tehran")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
			},
		},
		{
			name: "push again",
			fn: func() {
				err := repository.Push("France", "Paris")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
			},
		},
		{
			name: "pull",
			fn: func() {
				data, err := repository.Pull()
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				expectedData := models.NewData("Iran", "Tehran", 1)
				if data.String() != expectedData.String() {
					t.Errorf(errorMessageString, expectedData, data)
				}
			},
		},
		{
			name: "pull again",
			fn: func() {
				data, err := repository.Pull()
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				expectedData := models.NewData("France", "Paris", 2)
				if data.String() != expectedData.String() {
					t.Errorf(errorMessageString, expectedData, data)
				}
			},
		},
		{
			name: "push on empty queue",
			fn: func() {
				err := repository.Push("Germany", "Berlin")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
			},
		},
		{
			name: "pull another time",
			fn: func() {
				data, err := repository.Pull()
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				expectedData := models.NewData("Germany", "Berlin", 3)
				if data.String() != expectedData.String() {
					t.Errorf(errorMessageString, expectedData, data)
				}
			},
		},
		{
			name: "subscribe",
			fn: func() {
				err := repository.Subscribe(conn)
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
			},
		},
		{
			name: "push after subscribe",
			fn: func() {
				err := repository.Push("USA", "Washington")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}

				_, message, err := conn.ReadMessage()
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				expectedMessage := fmt.Sprintf(`["%s","%s"]`, "USA", "Washington")
				if string(message) != expectedMessage {
					t.Errorf(errorMessageString, expectedMessage, string(message))
				}
			},
		},
		{
			name: "subscribe again",
			fn: func() {
				err := repository.Subscribe(conn)
				if err != models.ErrSubscriberAlreadyExists {
					t.Errorf(errorMessage, models.ErrSubscriberAlreadyExists, err)
				}
			},
		},
		{
			name: "unsubscribe",
			fn: func() {
				err := repository.Unsubscribe(conn)
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
			},
		},
		{
			name: "unsubscribe again",
			fn: func() {
				err := repository.Unsubscribe(conn)
				if err != models.ErrSubscriberDoesNotExist {
					t.Errorf(errorMessage, models.ErrSubscriberDoesNotExist, err)
				}
			},
		},
		{
			name: "push after unsubscribe",
			fn: func() {
				err := repository.Push("Russia", "Moscow")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
			},
		},
		{
			name: "pull after unsubscribe",
			fn: func() {
				data, err := repository.Pull()
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				expectedData := models.NewData("Russia", "Moscow", 5)
				if data.String() != expectedData.String() {
					t.Errorf(errorMessageString, expectedData, data)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.fn()
		})
	}
}
