package v1

import (
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	"github.com/go-chi/chi/v5"
)

func NewRouter(router chi.Router, ms *service.MusicService, rs *service.RateService) {

	router.Route("/v1", func(r chi.Router) {
		NewMusicController(r, ms, rs)
	})

}
