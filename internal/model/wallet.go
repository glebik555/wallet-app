package model

type OperationRequest struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"`
	Amount        int64  `json:"amount"`
}

type BalanceResponse struct {
	WalletID string `json:"walletId"`
	Balance  int64  `json:"balance"`
}
