package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/teichx/go-expert-weather/internal/services"
)

type Handler struct {
	service *services.WeatherService
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ZipcodeRequest struct {
	CEP string `json:"cep"`
}

var ResponseInvalidZipcode = ErrorResponse{Message: "invalid zipcode"}
var ResponseCanNotFindZipcode = ErrorResponse{Message: "can not find zipcode"}
var ResponseInternalServerError = ErrorResponse{Message: "internal server error"}

func New(viaCEP services.ViaCEPClient, weatherAPI services.WeatherAPIClient) *Handler {
	service := services.New(viaCEP, weatherAPI)
	return &Handler{
		service: service,
	}
}

func (h *Handler) PostWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req ZipcodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseInvalidZipcode)
		return
	}

	zipcode := strings.TrimSpace(req.CEP)

	temp, err := h.service.GetWeatherByZipcode(r.Context(), zipcode)
	if err != nil {
		if errors.Is(err, services.ErrInvalidZipcode) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ResponseInvalidZipcode)
			return
		}

		if errors.Is(err, services.ErrCanNotFindZipcode) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ResponseCanNotFindZipcode)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ResponseInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temp)
}
