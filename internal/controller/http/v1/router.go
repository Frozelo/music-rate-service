package v1

import (
	mdl "github.com/Frozelo/music-rate-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(r chi.Router, userHandler *UserController, musicHandler *MusicController) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", SetupUserRoutes(userHandler))
		r.Group(func(r chi.Router) {
			r.Use(mdl.Auth)
			r.Mount("/music", SetupMusicRoutes(musicHandler))
		})
	})
}
