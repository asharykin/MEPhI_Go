package handlers

import (
	"banksystem/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
)

type CreditHandler struct {
	service *services.CreditService
}

func NewCreditHandler(service *services.CreditService) *CreditHandler {
	return &CreditHandler{service: service}
}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID int64   `json:"account_id"`
		Amount    float64 `json:"amount"`
		Term      int     `json:"term"` // срок в месяцах
		Rate      float64 `json:"rate"` // годовая процентная ставка
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(int64)
	credit, err := h.service.CreateCredit(r.Context(), userID, req.AccountID, req.Amount, req.Term, req.Rate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetUserCredits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)

	credits, err := h.service.GetUserCredits(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(credits)
}

func (h *CreditHandler) GetCredit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	credit, err := h.service.GetCredit(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if credit == nil {
		http.Error(w, "Credit not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	schedule, err := h.service.GetPaymentSchedule(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(schedule)
}
