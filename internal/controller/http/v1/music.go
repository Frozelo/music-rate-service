package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	"github.com/Frozelo/music-rate-service/internal/domain/usecase"
	music_usecase "github.com/Frozelo/music-rate-service/internal/domain/usecase/music"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type MusicController struct {
	mUcase    *music_usecase.MusicUsecase
	logger    logger.Interface
	validator *validator.Validate
}

func NewMusicController(mUcase *music_usecase.MusicUsecase, log logger.Interface) *MusicController {
	return &MusicController{
		mUcase:    mUcase,
		logger:    log,
		validator: validator.New(),
	}
}

func SetupMusicRoutes(musicHandler *MusicController) *chi.Mux {
	router := chi.NewRouter()
	{
		router.Get("/{musicId}/ratings", musicHandler.GetMusicRates)
		router.Get("/{musicId}/ratings/avg", musicHandler.GetMusicAverageRating)
		router.Post("/{musicId}/rate", musicHandler.RateMusic)
		router.Post("/{musicId}/nominate", musicHandler.NominateMusic)
	}
	return router
}

type ParamRateRequest struct {
	Param1 int `json:"p1" validate:"required,min=1,max=10"`
	Param2 int `json:"p2" validate:"required,min=1,max=10"`
	Param3 int `json:"p3" validate:"required,min=1,max=10"`
	Param4 int `json:"p4" validate:"required,min=1,max=10"`
}

type MusicRateRequest struct {
	Params  ParamRateRequest `json:"params" validate:"required"`
	Comment string           `json:"comment" validate:"required"`
}

type MusicNominateRequest struct {
	Nomination string `json:"nomination" validate:"required"`
}

func (mc *MusicController) RateMusic(w http.ResponseWriter, r *http.Request) {
	musicId, err := getMusicIdFromParam(r)
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid music ID: %w", err), mc.logger)
		return
	}

	var requestBody MusicRateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err), mc.logger)
		return
	}

	if err := mc.validator.Struct(requestBody.Params); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation failed: %w", err), mc.logger)
		return
	}

	if err := mc.mUcase.Rate(r.Context(), musicId, createRateDto(&requestBody)); err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to rate music: %w", err), mc.logger)
		return
	}

	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   map[string]string{"message": "ok"},
		Log:    mc.logger,
	})
}

func (mc *MusicController) NominateMusic(w http.ResponseWriter, r *http.Request) {
	musicId, err := getMusicIdFromParam(r)
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid music ID: %w", err), mc.logger)
		return
	}

	var requestBody MusicNominateRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err), mc.logger)
		return
	}

	log.Printf("Music ID: %d, Nomination: %s", musicId, requestBody.Nomination)

	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   map[string]string{"message": "ok"},
		Log:    mc.logger,
	})
}

func (mc *MusicController) DisplayMusics(w http.ResponseWriter, r *http.Request) {
	musics, err := mc.mUcase.GetAllMusic(r.Context())
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve musics: %w", err), mc.logger)
		return
	}

	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   musics,
		Log:    mc.logger,
	})
}

func (mc *MusicController) GetMusicRates(w http.ResponseWriter, r *http.Request) {
	musicId, err := getMusicIdFromParam(r)
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid music ID: %w", err), mc.logger)
		return
	}

	rates, err := mc.mUcase.GetAllMusicRates(r.Context(), musicId)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get music rates: %w", err), mc.logger)
		return
	}

	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   rates,
		Log:    mc.logger,
	})
}

func (mc *MusicController) GetMusicAverageRating(w http.ResponseWriter, r *http.Request) {
	musicId, err := strconv.Atoi(chi.URLParam(r, "musicId"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid music ID: %w", err), mc.logger)
		return
	}

	averageRating, err := mc.mUcase.GetAverageRating(r.Context(), musicId)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get average rating: %w", err), mc.logger)
		return
	}

	responseData := map[string]any{
		"musicId":   musicId,
		"avgRating": averageRating,
	}

	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   responseData,
		Log:    mc.logger,
	})
}

func getMusicIdFromParam(r *http.Request) (int, error) {
	return strconv.Atoi(chi.URLParam(r, "musicId"))
}

func createRateDto(req *MusicRateRequest) *usecase.MusicRateDto {
	return &usecase.MusicRateDto{
		Params: &entity.Rate{
			Param1: req.Params.Param1,
			Param2: req.Params.Param2,
			Param3: req.Params.Param3,
			Param4: req.Params.Param4,
		},
		Comment: req.Comment,
	}
}
