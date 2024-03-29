package internals

import (
	"io"
	"log"
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

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.RequestURI {
	case "/health":
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Println("Error writing response:", err)
		}
	case "/ws":
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading connection:", err)
			return
		}
		if _, ok := s.conns[conn]; ok {
			log.Println("Connection already exists:", conn.RemoteAddr())
			return
		}
		log.Println("New connection:", conn.RemoteAddr())

		s.conns[conn] = true
		s.ReadLoop(conn)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *Server) ReadLoop(conn *websocket.Conn) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by client:", conn.RemoteAddr())
				break
			}

			log.Println("Error reading message:", err)
			return
		}

		log.Println("Received message:", string(buf))

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
			response = "Usage: push <topic> <message> | pull | subscribe | unsubscribe"
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func (s *Server) Broadcast(message string) {
	for conn := range s.conns {
		go func(conn *websocket.Conn) {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Error writing message:", err)
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
