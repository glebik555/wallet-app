package server

import (
	"context"
	"testing"
	"wallet-app/internal/handlers"
	"wallet-app/internal/model"
)

type mockService struct{}

func (m *mockService) DoOperation(ctx context.Context, req model.OperationRequest) (model.BalanceResponse, error) {
	return model.BalanceResponse{WalletID: req.WalletID, Balance: 100}, nil
}

func (m *mockService) GetBalance(ctx context.Context, walletID string) (model.BalanceResponse, error) {
	return model.BalanceResponse{WalletID: walletID, Balance: 200}, nil
}

func TestNewServer(t *testing.T) {
	handler := handlers.NewWalletHandler(&mockService{})
	srv := New(handler, "8080")
	if srv.Addr != ":8080" {
		t.Errorf("expected :8080, got %s", srv.Addr)
	}
	if srv.Handler == nil {
		t.Error("expected non-nil handler")
	}
}
