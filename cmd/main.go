package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/internals"
)

func main() {
	server := internals.NewServer()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(server.Handler),
	}
	log.Println("Server started on :8080")

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Println("Error shutting down server:", err)
			return
		}

		log.Println("Shutting down server...")

		server.Close()
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalln("Error starting server:", err)
		}
	}
	<-idleConnsClosed
	log.Println("Server shutdown complete")
}
