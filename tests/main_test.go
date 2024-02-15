package tests

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/internals"
	"github.com/gorilla/websocket"
)

const errorMessage = `expected "%v", got "%v"`
const errorMessageString = `expected "%s", got "%s"`

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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

func TestMain(t *testing.T) {
	internalServer := internals.NewServer()
	defer internalServer.Close()

	server := httptest.NewServer(http.HandlerFunc(internalServer.Handler))
	defer server.Close()

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "health check",
			fn: func() {
				response, err := http.Get(server.URL + "/health")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				if response.StatusCode != http.StatusOK {
					t.Errorf(`expected status code %d, got %d`, http.StatusOK, response.StatusCode)
				} else {
					body, err := io.ReadAll(response.Body)
					if err != nil {
						t.Errorf(errorMessage, nil, err)
					} else if string(body) != "OK" {
						t.Errorf(`expected response body "OK", got "%s"`, string(body))
					}
				}
			},
		},
		{
			name: "not found",
			fn: func() {
				response, err := http.Get(server.URL + "/random-path")
				if err != nil {
					t.Errorf(errorMessage, nil, err)
				}
				if response.StatusCode != http.StatusNotFound {
					t.Errorf(`expected status code %d, got %d`, http.StatusNotFound, response.StatusCode)
				} else {
					body, err := io.ReadAll(response.Body)
					if err != nil {
						t.Errorf(errorMessage, nil, err)
					} else if string(body) != "Not found\n" {
						t.Errorf(`expected response body "Not found", got "%s"`, string(body))
					}
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
