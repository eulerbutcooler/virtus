package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/pkg/pagination"
	"github.com/eulerbutcooler/virtus/pkg/validator"
)

type envelope struct {
	Data  any        `json:"data,omitempty"`
	Error *errorBody `json:"error,omitempty"`
	Meta  any        `json:"meta,omitempty"`
}

type errorBody struct {
	Message string                     `json:"message"`
	Fields  validator.ValidationErrors `json:"fields,omitempty"`
}

// Writes a success envelope with the given status.
func OK(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, envelope{Data: data})
}

// Writes an error envelope with a plain message.
func Fail(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, envelope{Error: &errorBody{Message: msg}})
}

// Writes a 422 with per-field validation errors.
func FailValidation(w http.ResponseWriter, errs validator.ValidationErrors) {
	writeJSON(w, http.StatusUnprocessableEntity, envelope{
		Error: &errorBody{Message: "validation failed", Fields: errs},
	})
}

func List(w http.ResponseWriter, status int, data any, page pagination.Page) {
	writeJSON(w, status, envelope{Data: data, Meta: page})
}

// Maps a domain/service error to the right HTTP status.
func FromError(w http.ResponseWriter, err error) {
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		FailValidation(w, vErrs)
		return
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		Fail(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, domain.ErrConflict):
		Fail(w, http.StatusConflict, "resource already exists")
	case errors.Is(err, domain.ErrUnauthorized):
		Fail(w, http.StatusUnauthorized, "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		Fail(w, http.StatusForbidden, "forbidden")
	case errors.Is(err, domain.ErrInvalidInput):
		Fail(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, domain.ErrInsufficientFunds):
		Fail(w, http.StatusUnprocessableEntity, "insufficient pool funds")
	case errors.Is(err, domain.ErrInvalidState):
		Fail(w, http.StatusConflict, "invalid state transition")
	default:
		slog.Error("unhandled error", "error", err)
		Fail(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("response: encode failed", "error", err)
	}
}
