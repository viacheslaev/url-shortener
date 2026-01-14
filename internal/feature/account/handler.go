package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/viacheslaev/url-shortener/internal/server/httpx"
)

type RegisterHandler struct {
	service *AccountService
}

func NewAccountRegisterHandler(svc *AccountService) *RegisterHandler {
	return &RegisterHandler{
		service: svc,
	}
}

func (handler *RegisterHandler) RegisterAccount(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	publicId, err := handler.service.RegisterAccount(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	regResponse := createRegistrationResponse(publicId)
	httpx.WriteResponse(w, http.StatusCreated, regResponse)
}
