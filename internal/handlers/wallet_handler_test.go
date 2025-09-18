package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-app/internal/model"

	"github.com/stretchr/testify/assert"
)

type mockService struct{}

func (m *mockService) DoOperation(ctx context.Context, req model.OperationRequest) (model.BalanceResponse, error) {
	return model.BalanceResponse{WalletID: req.WalletID, Balance: 100}, nil
}

func (m *mockService) GetBalance(ctx context.Context, walletID string) (model.BalanceResponse, error) {
	return model.BalanceResponse{WalletID: walletID, Balance: 100}, nil
}

func TestHandleOperation_Deposit(t *testing.T) {
	ms := new(mockService)
	h := NewWalletHandler(ms)

	body := []byte(`{"walletId":"abc","operationType":"DEPOSIT","amount":100}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.HandleOperation(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.BalanceResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), resp.Balance)
}

func TestHandleOperation_InvalidJSON(t *testing.T) {
	ms := new(mockService)
	h := NewWalletHandler(ms)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader([]byte(`bad json`)))
	w := httptest.NewRecorder()

	h.HandleOperation(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
