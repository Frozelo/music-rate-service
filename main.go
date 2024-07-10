package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Frozelo/music-rate-service/pkg/httpserver"
)

func main() {
	httpServer := httpserver.New(httpserver.Port("8080"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Print("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		log.Fatal((fmt.Errorf("app - Run - httpServer.Notify: %w", err)))

	}
	err := httpServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
