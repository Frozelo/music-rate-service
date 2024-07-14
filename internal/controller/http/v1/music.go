package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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
		// h.POST(":musicId/nominate", r.NominateMusic)
	}
}

type MusicRateRequest struct {
	Param1 int `json:"p1" binding:"required,range1to10"`
	Param2 int `json:"p2" binding:"required,range1to10"`
	Param3 int `json:"p3" binding:"required,range1to10"`
	Param4 int `json:"p4" binding:"required,range1to10"`
}

type MusicNominateRequest struct {
	Nomination string `json:"nomination" binding:"required"`
}

// Custom validator for range 1 to 10
func range1to10(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	return value >= 1 && value <= 10
}

func (mc *MusicController) RateMusic(c *gin.Context) {
	ctx := c.Request.Context()
	musicId, err := mc.getMusicIdFromParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid music ID: " + err.Error()})
		return
	}

	var requestBody MusicRateRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	rateDto := mc.createRateDto(&requestBody)
	rate := mc.rs.CalculateRate(rateDto)

	if err := mc.ms.Rate(ctx, musicId, rate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rate music: " + err.Error()})
		return
	}

	log.Printf("Successfully rated music with ID %d: %+v", musicId, rateDto)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// func (mc *MusicController) NominateMusic(c *gin.Context) {
// 	musicId, err := mc.getMusicIdFromParam(c)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid music ID: " + err.Error()})
// 		return
// 	}

// 	var requestBody MusicNominateRequest
// 	if err := c.ShouldBindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
// 		return
// 	}

// 	if err := mc.ms.Nominate(ctx, musicId, requestBody.Nomination); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to nominate music: " + err.Error()})
// 		return
// 	}

// 	log.Printf("Successfully nominated music with ID %d as %s", musicId, requestBody.Nomination)
// 	c.JSON(http.StatusOK, gin.H{"message": "ok"})
// }

func (mc *MusicController) getMusicIdFromParam(c *gin.Context) (int, error) {
	musicId, err := strconv.Atoi(c.Param("musicId"))
	if err != nil {
		return 0, err
	}
	return musicId, nil
}

func (mc *MusicController) createRateDto(request *MusicRateRequest) *entity.Rate {
	return &entity.Rate{
		Param1: request.Param1,
		Param2: request.Param2,
		Param3: request.Param3,
		Param4: request.Param4,
	}
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("range1to10", range1to10)
	}
}
