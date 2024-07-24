package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Frozelo/music-rate-service/internal/domain/entity"
	user_usecase "github.com/Frozelo/music-rate-service/internal/domain/usecase/user"
	"github.com/Frozelo/music-rate-service/pkg/httpserver"
	jwt_service "github.com/Frozelo/music-rate-service/pkg/jwt"
	"github.com/Frozelo/music-rate-service/pkg/logger"
	"github.com/Frozelo/music-rate-service/pkg/oauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	uUcase *user_usecase.UserUsecase
	logger logger.Interface // Добавлен логгер для логирования ошибок и информации
}

func NewUserController(uUcase *user_usecase.UserUsecase, log logger.Interface) *UserController {
	return &UserController{uUcase: uUcase, logger: log}
}

func SetupUserRoutes(userHandler *UserController) chi.Router {
	router := chi.NewRouter()
	router.Post("/register", userHandler.CreateUser)
	router.Get("/auth/github/login", oauth.HandleGitHubLogin)
	router.Get("/auth/github/callback", oauth.HandleGitHubCallback)
	router.Post("/auth/login", userHandler.Login)
	return router
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func init() {
	validate = validator.New()
}

func (req *CreateUserRequest) Bind() error {
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
	var requestBody CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	if err := requestBody.Bind(); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	user := user_usecase.CreateUserDto{
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	if err := uc.uUcase.RegisterUser(r.Context(), user); err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}
	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusCreated,
		Data:   map[string]string{"message": "user created"},
		Log:    uc.logger,
	})
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var requestBody LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	user := user_usecase.LoginUserDto{
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	authenticatedUser, err := uc.uUcase.LoginUser(ctx, user)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}

	token, err := jwt_service.GenerateJWT(authenticatedUser)
	if err != nil {
		httpserver.WriteError(w, http.StatusInternalServerError, err, uc.logger)
		return
	}

	response := map[string]string{
		"token": token,
	}
	httpserver.WriteJSONResponse(w, httpserver.ResponseConfig{
		Status: http.StatusOK,
		Data:   response,
		Log:    uc.logger,
	})
}

func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	if err != nil {
		httpserver.WriteError(w, http.StatusBadRequest, err, uc.logger)
		return
	}

	var requestBody CreateUserRequest
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
