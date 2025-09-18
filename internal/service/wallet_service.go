package service

import (
	"context"
	"errors"
	"time"

	"wallet-app/internal/model"
	"wallet-app/internal/repo"
)

var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrInvalidOperationType = errors.New("operationType must be DEPOSIT or WITHDRAW")
	ErrInvalidAmount        = errors.New("amount must be > 0")
)

type WalletService interface {
	DoOperation(ctx context.Context, req model.OperationRequest) (model.BalanceResponse, error)
	GetBalance(ctx context.Context, walletID string) (model.BalanceResponse, error)
}

type walletService struct {
	repo *repo.WalletRepo
}

func NewWalletService(r *repo.WalletRepo) WalletService {
	return &walletService{repo: r}
}

func (s *walletService) DoOperation(ctx context.Context, req model.OperationRequest) (model.BalanceResponse, error) {
	if req.Amount <= 0 {
		return model.BalanceResponse{}, ErrInvalidAmount
	}
	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		return model.BalanceResponse{}, ErrInvalidOperationType
	}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	resp, err := s.repo.DoOperation(ctx, req)
	if err != nil {
		if errors.Is(err, repo.ErrInsufficientFunds) {
			return model.BalanceResponse{}, ErrInsufficientFunds
		}
		return model.BalanceResponse{}, err
	}
	return resp, nil
}

func (s *walletService) GetBalance(ctx context.Context, walletID string) (model.BalanceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	return s.repo.GetBalance(ctx, walletID)
}
