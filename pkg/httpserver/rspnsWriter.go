package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Frozelo/music-rate-service/pkg/logger"
)

// HTTPError представляет структуру для ошибок в формате JSON.
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Предопределенные сообщения об ошибках
const (
	ErrorMessageNotFound       = "Resource not found"
	ErrorMessageInternalServer = "Internal Server Error"
	ErrorMessageBadRequest     = "Bad Request"
	ErrorMessageUnauthorized   = "Unauthorized"
	ErrorMessageForbidden      = "Forbidden"
)

type ResponseConfig struct {
	Status  int
	Data    any
	Log     logger.Interface
	Headers map[string]string
	Err     error
}

func WriteJSONResponse(w http.ResponseWriter, config ResponseConfig) {
	startTime := time.Now()

	// Установка заголовков
	w.Header().Set("Content-Type", "application/json")
	for key, value := range config.Headers {
		w.Header().Set(key, value)
	}

	w.WriteHeader(config.Status)
	if err := json.NewEncoder(w).Encode(config.Data); err != nil {
		config.Log.Error("Failed to encode JSON response: " + err.Error())
		http.Error(w, ErrorMessageInternalServer, http.StatusInternalServerError)
		return
	}

	duration := time.Since(startTime)
	logResponse(config.Status, config.Data, duration, config.Log)
}

// WriteError отправляет ошибку в формате JSON с предопределенным сообщением.
func WriteError(w http.ResponseWriter, code int, err error, log logger.Interface) {
	message := map[int]string{
		http.StatusNotFound:            ErrorMessageNotFound,
		http.StatusInternalServerError: ErrorMessageInternalServer,
		http.StatusBadRequest:          ErrorMessageBadRequest,
		http.StatusUnauthorized:        ErrorMessageUnauthorized,
		http.StatusForbidden:           ErrorMessageForbidden,
	}[code]

	if message == "" {
		message = "An unexpected error occurred"
	}

	if err != nil {
		message += ": " + err.Error()
	}

	WriteJSONResponse(w, ResponseConfig{
		Status: code,
		Data:   HTTPError{Code: code, Message: message},
		Log:    log,
	})
}

// logResponse логирует информацию о статусе ответа, данных и времени выполнения.
func logResponse(status int, data any, duration time.Duration, log logger.Interface) {
	log.Info("Response: status=%d, data=%+v, duration=%s", status, data, duration)
}
