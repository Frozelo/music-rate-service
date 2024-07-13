package v1

import (
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.RouterGroup, ms *service.MusicService, rs *service.RateService) {

	h := handler.Group("v1")
	{
		NewMusicController(h, ms, rs)
	}
}
