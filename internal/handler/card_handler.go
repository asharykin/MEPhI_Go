package handler

import (
	"banksystem/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
)

type CardHandler struct {
	service *service.CardService
}

func NewCardHandler(service *service.CardService) *CardHandler {
	return &CardHandler{service: service}
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var request struct {
		AccountID int64 `json:"account_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	card, err := h.service.CreateCard(request.AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) GetUserCards(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	cards, err := h.service.GetUserCards(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cards)
}

func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	cardID, err := strconv.ParseInt(r.URL.Query().Get("card_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	card, err := h.service.GetCard(cardID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if card == nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(card)
}
