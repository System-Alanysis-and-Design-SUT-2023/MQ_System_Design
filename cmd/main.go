package main

import (
	"net/http"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/tools"
)

func main() {
	server := tools.NewServer()
	http.Handle("/ws", http.HandlerFunc(server.HandleWS))
	http.ListenAndServe(":8080", nil)
}
