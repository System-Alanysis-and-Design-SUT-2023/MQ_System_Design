package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/internals"
	"github.com/System-Alanysis-and-Design-SUT-2023/MQ_System_Design/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalln("Error loading config:", err)
	}
	listenPort := config.ListenPort

	server := internals.NewServer()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: http.HandlerFunc(server.Handler),
	}
	log.Printf("Listening on port :%d...\n", listenPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalln("Error starting server:", err)
			}
		}
	}()

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
	log.Println("Server shutdown complete!")
}
