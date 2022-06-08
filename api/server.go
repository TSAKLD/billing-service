package api

import (
	"avitoTZ/service"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	port   string
	router *mux.Router
}

func NewServer(serv *service.UserService, p string) *Server {
	hand := NewHandler(serv)
	r := mux.NewRouter()

	r.HandleFunc("/users", hand.AddUser).Methods(http.MethodPost)
	r.HandleFunc("/users/balance", hand.UserBalance).Methods(http.MethodGet)
	r.HandleFunc("/users/balance", hand.ChangeUserBalance).Methods(http.MethodPatch)
	r.HandleFunc("/users/balance-transfer", hand.TransferMoney).Methods(http.MethodPatch)
	r.HandleFunc("/users/currency-balance", hand.BalanceInCurrency).Methods(http.MethodGet)
	r.HandleFunc("/users/transactions", hand.TransactionList).Methods(http.MethodGet)

	return &Server{
		port:   p,
		router: r,
	}
}

func (r *Server) ListenAndServe() error {
	fmt.Println("Server is listening...")

	err := http.ListenAndServe(":"+r.port, r.router)
	return err
}
