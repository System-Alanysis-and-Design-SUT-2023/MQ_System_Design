package tools

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
	repo  *Repository
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection: ", err)
		return
	}
	if _, ok := s.conns[conn]; ok {
		fmt.Println("Connection already exists: ", conn.RemoteAddr())
		return
	}

	fmt.Println("New connection from client: ", conn.RemoteAddr())
	s.conns[conn] = true
	s.ReadLoop(conn)
}

func (s *Server) ReadLoop(conn *websocket.Conn) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client: ", conn.RemoteAddr())
				break
			}

			fmt.Println("Error reading message: ", err)
			continue
		}

		if len(buf) == 0 {
			continue
		}

		fmt.Println("Received message: ", string(buf))

		response := ""
		commands := strings.Fields(string(buf))
		switch commands[0] {
		case "push":
			response = HandlePush(s.repo, commands[1:]...)
		case "pull":
			response = HandlePull(s.repo)
		case "subscribe":
			response = HandleSubscribe(s.repo, conn)
		case "unsubscribe":
			response = HandleUnsubscribe(s.repo, conn)
		default:
			response = "Invalid command!\nUsage: push <topic> <message> | pull | subscribe | unsubscribe"
		}

		conn.WriteMessage(websocket.TextMessage, []byte(response))
	}
}

func (s *Server) Broadcast(message string) {
	for conn := range s.conns {
		go func(conn *websocket.Conn) {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				fmt.Println("Error writing message: ", err)
				return
			}
		}(conn)
	}
}

func (s *Server) Close() {
	for conn := range s.conns {
		delete(s.conns, conn)
		conn.Close()
	}
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		repo:  NewRepository(),
	}
}
