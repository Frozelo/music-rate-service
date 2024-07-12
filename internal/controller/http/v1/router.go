package controller

import (
	"net/http"
	"strconv"

	"github.com/Frozelo/music-rate-service/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type MusicController struct {
	musicService *service.MusicService
}

func NewMusicController(musicService *service.MusicService) *MusicController {
	return &MusicController{
		musicService: musicService,
	}
}

func (mc *MusicController) RateMusic(c *gin.Context) {
	musicId, err := strconv.Atoi(c.Param("musicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid music ID"})
		return
	}

	var request struct {
		Rate int `json:"rate" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = mc.musicService.Rate(musicId, request.Rate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Music rated successfully"})
}
