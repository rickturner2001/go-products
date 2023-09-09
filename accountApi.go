package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == http.MethodGet {
		return s.handleGetAccount(w, r)
	}

	if r.Method == http.MethodPost {

		return s.handleCreateAccount(w, r)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteProduct(w, r)
	}

	if r.Method == http.MethodPatch {
		return s.handleUpdateProduct(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accountId := mux.Vars(r)["id"]

	

	if accountId == ""{
		accounts, err := s.store.GetAllAccounts()

		if err != nil {
			return WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
		}

		return WriteJSON(w, http.StatusOK, accounts)
	}

	parsedId, err := strconv.ParseInt(accountId, 0, 8)

	if err != nil {
		return err
	}

	account , err := s.store.GetAccountByID(int(parsedId))

	if err != nil {
		return err
	}

	if account == nil {
		return WriteJSON(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Could not find account with id: %d", parsedId)})
	}

	return WriteJSON(w, http.StatusOK, account)

	
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	account, err := extractAccountFromRequest(r)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "Could not unmarshal account data"})
	}

	tokenString, err := CreateJWT(account)

	if err != nil {
		return err
	}

	fmt.Println("Token string: ", tokenString)

	err = s.store.CreateAccount(account)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func extractAccountFromRequest(r *http.Request) (*Account, error) {
	createProductReq := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(createProductReq); err != nil {
		return nil, err
	}

	account := NewAccount(createProductReq.Username)
	return account, nil

}
