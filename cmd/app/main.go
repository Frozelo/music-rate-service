package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Frozelo/music-rate-service/config"
	v1 "github.com/Frozelo/music-rate-service/internal/controller/http/v1"
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	music_usecase "github.com/Frozelo/music-rate-service/internal/domain/usecase/music"
	user_usecase "github.com/Frozelo/music-rate-service/internal/domain/usecase/user"
	postgres_repository "github.com/Frozelo/music-rate-service/internal/repository/postgres"
	"github.com/Frozelo/music-rate-service/internal/storage"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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

	musicRepo := postgres_repository.NewMusicRepository(storage.Conn)
	musicService := service.NewMusicService(musicRepo)
	rateService := service.NewRateService()
	musicUsecase := music_usecase.NewMusicUsecase(musicService, rateService)

	userRepo := postgres_repository.NewUserRepository(storage.Conn)
	userService := service.NewUserService(userRepo)
	userUsecase := user_usecase.NewUserUsecase(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	apiGroup := chi.NewRouter()
	v1.NewRouter(apiGroup, musicUsecase, userUsecase, l)
	r.Mount("/api", apiGroup)

	l.Info("starting new http server")
	httpServer := httpserver.New(r, httpserver.Port(cfg.Server.Port))
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
