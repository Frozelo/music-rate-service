package v1

import (
	"log"
	"net/http"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type MusicController struct {
	ms *service.MusicService
	rs *service.RateService
}

func NewMusicController(handler *gin.RouterGroup, ms *service.MusicService, rs *service.RateService) {
	r := &MusicController{ms: ms, rs: rs}
	h := handler.Group("/music")
	{
		h.POST(":musicId/rate", r.RateMusic)
	}
}

func (mc *MusicController) RateMusic(c *gin.Context) {

	var request struct {
		Param1 int `json:"p1" binding:"required"`
		Param2 int `json:"p2" binding:"required"`
		Param3 int `json:"p3" binding:"required"`
		Param4 int `json:"p4" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Print(request)
	rateDto := &entity.Rate{
		Param1: request.Param1,
		Param2: request.Param2,
		Param3: request.Param3,
		Param4: request.Param4,
	}
	rate := mc.rs.CalculateRate(rateDto)

	c.JSON(http.StatusOK, gin.H{"message": rate})
}
