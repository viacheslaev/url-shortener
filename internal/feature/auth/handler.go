package auth

import (
	"encoding/json"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(authService *AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	token, err := handler.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpx.WriteErr(w, http.StatusUnauthorized, ErrInvalidCredentials.Error())
		return
	}

	httpx.WriteResponse(w, http.StatusOK, LoginResponse{Token: token})
}
