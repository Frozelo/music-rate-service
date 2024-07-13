package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/Frozelo/music-rate-service/internal/controller/http/v1"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	memmory_repository "github.com/Frozelo/music-rate-service/internal/repository/memmory"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	l := logger.New("debug")

	// connString := "-"
	// storage, err := storage.New(connString)

	// if err != nil {
	// 	l.Fatal("Unable to connect to database: %v\n", err)
	// }
	// defer storage.Close()

	musicRepo := memmory_repository.NewMusicRepository()
	music := &entity.Music{Name: "Song A", Author: "Author A", Rate: 5}
	musicRepo.Create(music)

	musicService := service.NewMusicService(musicRepo)
	rateService := service.NewRateService()

	handler := gin.New()
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	apiGroup := handler.Group("/api")
	{
		v1.NewRouter(apiGroup, musicService, rateService)
	}

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
