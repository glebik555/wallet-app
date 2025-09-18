package repo

import (
	"context"
	"errors"

	"wallet-app/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type WalletRepo struct {
	db *pgxpool.Pool
}

func NewWalletRepo(db *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) GetBalance(ctx context.Context, walletID string) (model.BalanceResponse, error) {
	row := r.db.QueryRow(ctx, `SELECT uuid, balance FROM wallets WHERE uuid=$1`, walletID)

	var resp model.BalanceResponse
	if err := row.Scan(&resp.WalletID, &resp.Balance); err != nil {
		return model.BalanceResponse{}, err
	}
	return resp, nil
}

func (r *WalletRepo) DoOperation(ctx context.Context, req model.OperationRequest) (model.BalanceResponse, error) {
	switch req.OperationType {

	case "DEPOSIT":
		tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return model.BalanceResponse{}, err
		}
		defer tx.Rollback(ctx)

		var balance int64
		err = tx.QueryRow(ctx,
			`SELECT balance FROM wallets WHERE uuid=$1 FOR UPDATE`,
			req.WalletID).Scan(&balance)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				err = tx.QueryRow(ctx,
					`INSERT INTO wallets (uuid, balance)
					 VALUES ($1, $2)
					 RETURNING balance`,
					req.WalletID, req.Amount).Scan(&balance)
				if err != nil {
					return model.BalanceResponse{}, err
				}
			} else {
				return model.BalanceResponse{}, err
			}
		} else {
			err = tx.QueryRow(ctx,
				`UPDATE wallets
				 SET balance = balance + $1, updated_at = now()
				 WHERE uuid=$2
				 RETURNING balance`,
				req.Amount, req.WalletID).Scan(&balance)
			if err != nil {
				return model.BalanceResponse{}, err
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return model.BalanceResponse{}, err
		}
		return model.BalanceResponse{WalletID: req.WalletID, Balance: balance}, nil

	case "WITHDRAW":
		tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
		if err != nil {
			return model.BalanceResponse{}, err
		}
		defer tx.Rollback(ctx)

		var balance int64
		err = tx.QueryRow(ctx,
			`SELECT balance FROM wallets WHERE uuid=$1 FOR UPDATE`,
			req.WalletID).Scan(&balance)
		if err != nil {
			return model.BalanceResponse{}, err
		}

		if balance < req.Amount {
			return model.BalanceResponse{}, ErrInsufficientFunds
		}

		var newBalance int64
		err = tx.QueryRow(ctx,
			`UPDATE wallets
			 SET balance = balance - $1, updated_at = now()
			 WHERE uuid=$2
			 RETURNING balance`,
			req.Amount, req.WalletID).Scan(&newBalance)
		if err != nil {
			return model.BalanceResponse{}, err
		}

		if err := tx.Commit(ctx); err != nil {
			return model.BalanceResponse{}, err
		}
		return model.BalanceResponse{WalletID: req.WalletID, Balance: newBalance}, nil

	default:
		return model.BalanceResponse{}, errors.New("invalid operation type")
	}
}
