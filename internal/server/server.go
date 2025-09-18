package server

import (
	"net/http"
	"time"

	"wallet-app/internal/handlers"

	"github.com/gorilla/mux"
)

var ErrServerClosed = http.ErrServerClosed

func New(handler *handlers.WalletHandler, port string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallet", handler.HandleOperation).Methods("POST")
	r.HandleFunc("/api/v1/wallets/{uuid}", handler.HandleGetBalance).Methods("GET")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  300 * time.Second,
	}
	return srv
}
