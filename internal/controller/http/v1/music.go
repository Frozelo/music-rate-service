package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type MusicController struct {
	ms *service.MusicService
	rs *service.RateService
}

func NewMusicController(router chi.Router, ms *service.MusicService, rs *service.RateService) {
	mc := &MusicController{ms: ms, rs: rs}
	router.Post("/music/{musicId}/rate", mc.RateMusic)
	// router.Post("/music/{musicId}/nominate", mc.NominateMusic)
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

var validate *validator.Validate

// Custom validator for range 1 to 10
func range1to10(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	return value >= 1 && value <= 10
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("range1to10", range1to10)
}

func (req *MusicRateRequest) Bind(r *http.Request) error {
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func (mc *MusicController) RateMusic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	musicId, err := mc.getMusicIdFromParam(r)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid music ID: " + err.Error()})
		return
	}

	var requestBody MusicRateRequest
	if err := render.Bind(r, &requestBody); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request data: " + err.Error()})
		return
	}

	rateDto := mc.createRateDto(&requestBody)
	rate := mc.rs.CalculateRate(rateDto)

	if err := mc.ms.Rate(ctx, musicId, rate); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to rate music: " + err.Error()})
		return
	}

	log.Printf("Successfully rated music with ID %d: %+v", musicId, rateDto)
	render.JSON(w, r, map[string]string{"message": "ok"})
}

// func (mc *MusicController) NominateMusic(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	musicId, err := mc.getMusicIdFromParam(r)
// 	if err != nil {
// 		render.Status(r, http.StatusBadRequest)
// 		render.JSON(w, r, map[string]string{"error": "Invalid music ID: " + err.Error()})
// 		return
// 	}

// 	var requestBody MusicNominateRequest
// 	if err := render.Bind(r, &requestBody); err != nil {
// 		render.Status(r, http.StatusBadRequest)
// 		render.JSON(w, r, map[string]string{"error": "Invalid request data: " + err.Error()})
// 		return
// 	}

// 	if err := mc.ms.Nominate(ctx, musicId, requestBody.Nomination); err != nil {
// 		render.Status(r, http.StatusInternalServerError)
// 		render.JSON(w, r, map[string]string{"error": "Failed to nominate music: " + err.Error()})
// 		return
// 	}

// 	log.Printf("Successfully nominated music with ID %d as %s", musicId, requestBody.Nomination)
// 	render.JSON(w, r, map[string]string{"message": "ok"})
// }

func (mc *MusicController) getMusicIdFromParam(r *http.Request) (int, error) {
	musicId, err := strconv.Atoi(chi.URLParam(r, "musicId"))
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
