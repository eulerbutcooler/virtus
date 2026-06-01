package v1

import (
	"net/http"

	"github.com/eulerbutcooler/virtus/internal/domain"
	"github.com/eulerbutcooler/virtus/internal/handler/response"
	"github.com/eulerbutcooler/virtus/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type authResponse struct {
	User  *domain.User       `json:"user"`
	Token *service.TokenPair `json:"token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := response.DecodeAndValidate(r, &req); err != nil {
		response.FromError(w, err)
		return
	}
	user, tokens, err := h.auth.Register(r.Context(), service.RegisterInput{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusCreated, authResponse{User: user, Token: tokens})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := response.DecodeAndValidate(r, &req); err != nil {
		response.FromError(w, err)
		return
	}
	user, tokens, err := h.auth.Login(r.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, authResponse{User: user, Token: tokens})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := response.DecodeAndValidate(r, &req); err != nil {
		response.FromError(w, err)
		return
	}
	tokens, err := h.auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, tokens)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := response.DecodeAndValidate(r, &req); err != nil {
		response.FromError(w, err)
		return
	}
	if err := h.auth.Logout(r.Context(), req.RefreshToken); err != nil {
		response.FromError(w, err)
		return
	}
	response.OK(w, http.StatusOK, map[string]string{"status": "logged out"})
}
