package v1

import (
	music_usecase "github.com/Frozelo/music-rate-service/internal/domain/usecase/music"
	"github.com/go-chi/chi/v5"
)

func NewRouter(router chi.Router, mUcase *music_usecase.MusicUsecase) {

	router.Route("/v1", func(r chi.Router) {
		NewMusicController(r, mUcase)
	})

}
