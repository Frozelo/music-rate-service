package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/Frozelo/music-rate-service/internal/conrtoller/http/v1"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	l := logger.New("debug")

	handler := gin.New()
	v1.NewRouter(handler)

	l.Info("starting new http server")
	httpServer := httpserver.New(handler, httpserver.Port("8080"))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error((fmt.Errorf("app - Run - httpServer.Notify: %w", err)))

	}
	err := httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
