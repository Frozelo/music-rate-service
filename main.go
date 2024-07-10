package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/Frozelo/music-rate-service/internal/conrtoller/http/v1"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/gin-gonic/gin"
)

func main() {
	handler := gin.New()
	v1.NewRouter(handler)

	httpServer := httpserver.New(handler, httpserver.Port("8080"))

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
