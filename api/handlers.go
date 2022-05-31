package api

import (
	"avitoTZ/entity"
	"avitoTZ/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Handler struct {
	us *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	h := Handler{
		us: userService,
	}

	return &h
}

func (hdr Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var u entity.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	err = u.Validate()
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	u, err = hdr.us.AddUser(u)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (hdr Handler) UserBalance(w http.ResponseWriter, r *http.Request) {
	var u entity.User

	QID := r.URL.Query().Get("id")

	id, err := strconv.Atoi(QID)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	u, err = hdr.us.UserBalance(id)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(u)
}

func (hdr Handler) ChangeUserBalance(w http.ResponseWriter, r *http.Request) {
	var tr entity.Transaction

	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	if tr.Retriever == nil {
		err := errors.New("input retriever ID")
		sendError(w, err, http.StatusBadRequest)
		return
	}

	if tr.Amount == nil {
		err := errors.New("input operation amount")
		sendError(w, err, http.StatusBadRequest)
		return
	}

	err = hdr.us.ChangeUserBalance(tr)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("Operation Completed")
}

func (hdr Handler) TransferMoney(w http.ResponseWriter, r *http.Request) {
	var tr entity.Transaction

	err := json.NewDecoder(r.Body).Decode(&tr)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	if tr.Amount == nil || tr.Sender == nil || tr.Retriever == nil {
		err := errors.New("input all the required fields")
		sendError(w, err, http.StatusBadRequest)
		return
	}

	if *tr.Amount <= 0 {
		err := errors.New("invalid amount")
		sendError(w, err, http.StatusBadRequest)
		return
	}

	err = hdr.us.TransferMoney(tr)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode("Transaction Completed")
}

func (hdr Handler) BalanceInCurrency(w http.ResponseWriter, r *http.Request) {
	QID := r.URL.Query().Get("id")

	id, err := strconv.Atoi(QID)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	c := r.URL.Query().Get("currency")

	result, err := hdr.us.BalanceInCurrency(id, c)
	if err != nil {
		sendError(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(result)
}
