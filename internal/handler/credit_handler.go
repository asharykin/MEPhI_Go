package handler

import (
	"banksystem/internal/dto"
	"banksystem/internal/logger"
	"banksystem/internal/middleware"
	"banksystem/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type CreditHandler struct {
	service *service.CreditService
}

func NewCreditHandler(service *service.CreditService) *CreditHandler {
	return &CreditHandler{service: service}
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateCredit(r.Context(), userID, &req)
	if err != nil {
		logger.Error("Failed to create credit via handler", "error", err, "user_id", userID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	creditID := vars["id"]

	if creditID == "" {
		http.Error(w, "Credit ID is required", http.StatusBadRequest)
		return
	}

	schedule, err := h.service.GetCreditSchedule(r.Context(), creditID, userID)
	if err != nil {
		logger.Error("Failed to get credit schedule via handler", "error", err, "credit_id", creditID, "user_id", userID)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}
