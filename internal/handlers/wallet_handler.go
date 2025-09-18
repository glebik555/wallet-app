package handlers

import (
	"encoding/json"
	"net/http"

	"wallet-app/internal/model"
	"wallet-app/internal/service"

	"github.com/gorilla/mux"
)

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(s service.WalletService) *WalletHandler {
	return &WalletHandler{service: s}
}

// POST /api/v1/wallet
func (h *WalletHandler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req model.OperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if req.WalletID == "" {
		http.Error(w, `{"error":"missing walletId"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.service.DoOperation(r.Context(), req)
	if err != nil {
		switch err {
		case service.ErrInsufficientFunds:
			http.Error(w, `{"error":"insufficient funds"}`, http.StatusBadRequest)
			return
		case service.ErrInvalidOperationType, service.ErrInvalidAmount:
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		default:
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/wallets/{uuid}
func (h *WalletHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	walletID := vars["uuid"]
	if walletID == "" {
		http.Error(w, `{"error":"missing uuid"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetBalance(r.Context(), walletID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
