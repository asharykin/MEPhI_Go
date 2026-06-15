package handler

import (
	"banksystem/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type CreditPaymentHandler struct {
	service *services.CreditPaymentService
}

func NewCreditPaymentHandler(service *services.CreditPaymentService) *CreditPaymentHandler {
	return &CreditPaymentHandler{service: service}
}

func (h *CreditPaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CreditID int64     `json:"credit_id"`
		Amount   float64   `json:"amount"`
		DueDate  time.Time `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.DueDate.Before(time.Now()) {
		http.Error(w, "Due date cannot be in the past", http.StatusBadRequest)
		return
	}

	payment, err := h.service.CreatePayment(r.Context(), req.CreditID, req.Amount, req.DueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func (h *CreditPaymentHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	paymentID, err := strconv.ParseInt(r.URL.Query().Get("payment_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	err = h.service.ProcessPayment(r.Context(), paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CreditPaymentHandler) GetPaymentsByCreditID(w http.ResponseWriter, r *http.Request) {
	creditID, err := strconv.ParseInt(r.URL.Query().Get("credit_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid credit ID", http.StatusBadRequest)
		return
	}

	payments, err := h.service.GetPaymentsByCreditID(r.Context(), creditID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payments)
}

func (h *CreditPaymentHandler) GetPendingPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := h.service.GetPendingPayments(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payments)
}
