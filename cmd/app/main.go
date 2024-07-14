package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Frozelo/music-rate-service/config"
	v1 "github.com/Frozelo/music-rate-service/internal/controller/http/v1"
	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	memmory_repository "github.com/Frozelo/music-rate-service/internal/repository/memmory"
	"github.com/Frozelo/music-rate-service/internal/storage"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/gin-gonic/gin"
)

const configPath string = "config/config.yml"

func main() {
	log.Print("Config initialzation")
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Config initialization error: %s", err)
	}
	log.Print("Successful config initialization")

	l := logger.New(cfg.Log.Level)
	l.Info("Successful logger initialization")

	storage, err := storage.New(cfg.Database.ConnString)

	if err != nil {
		l.Fatal("Unable to connect to database: %v\n", err)
	}
	defer storage.Close()

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
	httpServer := httpserver.New(handler, httpserver.Port(cfg.Server.Port))
	l.Info("Successful server startup on port %s", cfg.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		l.Error((fmt.Errorf("app - Run - httpServer.Notify: %w", err)))

	}
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
