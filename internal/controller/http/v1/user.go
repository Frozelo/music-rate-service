package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	user_usecase "github.com/Frozelo/music-rate-service/internal/domain/usecase/user"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	uUcase *user_usecase.UserUsecase
	logger logger.Interface // Добавлен логгер для логирования ошибок и информации
}

func NewUserController(router chi.Router, uUcase *user_usecase.UserUsecase, log logger.Interface) {
	uc := &UserController{uUcase: uUcase, logger: log}
	router.Get("/users", uc.GetAllUsers)
	router.Get("/users/{userId}", uc.GetUserByID)
	router.Post("/users", uc.CreateUser)
	router.Put("/users/{userId}", uc.UpdateUser)
	router.Delete("/users/{userId}", uc.DeleteUser)
}

type UserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func init() {
	validate = validator.New()
}

func (req *UserRequest) Bind() error {
	return validate.Struct(req)
}

func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := uc.uUcase.GetAllUsers(ctx)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}
	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   users,
		Log:    uc.logger,
	})
}

func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	user, err := uc.uUcase.GetUserByID(ctx, userId)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}

	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   user,
		Log:    uc.logger,
	})
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	if err := requestBody.Bind(); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	user := &entity.User{
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	if err := uc.uUcase.CreateUser(r.Context(), user); err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}
	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusCreated,
		Data:   map[string]string{"message": "user created"},
		Log:    uc.logger,
	})
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	var requestBody UserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	if err := requestBody.Bind(); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	user := &entity.User{
		ID:       userId,
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}

	if err := uc.uUcase.UpdateUser(r.Context(), user); err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}
	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   map[string]string{"message": "user updated"},
		Log:    uc.logger,
	})
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	if err := uc.uUcase.DeleteUser(r.Context(), userId); err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}
	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   map[string]string{"message": "user deleted"},
		Log:    uc.logger,
	})
}
