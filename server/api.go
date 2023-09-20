package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/product", makeHTTPHandleFunc(s.HandleProduct))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))

	router.HandleFunc("/product/{id}", withJWTAuth(makeHTTPHandleFunc(s.HandleProduct), s.store))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.HandleAccount), s.store))

	router.HandleFunc("/status", makeHTTPHandleFunc(statusHandle))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatalf("Could not listen  on port %s: %v", s.listenAddr, err)
	}
}

type StatusResponse struct {
	message string
	success bool
}

func getStatusResponse() *StatusResponse {
	return &StatusResponse{
		message: "Service is working",
		success: true,
	}
}

func statusHandle(rw http.ResponseWriter, r *http.Request) error {
	err := WriteJSON(rw, http.StatusOK, getStatusResponse())
	if err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if err != nil {
				log.Printf("could not write json response: %v", err)
				return
			}
		}
	}
}

func CreateJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":         15000,
		"accountIdentifier": account.AccountIdentifier,
	}

	secret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func accessDenied(w http.ResponseWriter) {
	err := WriteJSON(w, http.StatusForbidden, ApiError{Error: "Access denied"})
	if err != nil {
		return
	}
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)

		if err != nil {
			accessDenied(w)
			return
		}

		if !token.Valid {
			accessDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		accountId, err := ExtractIdFromRequest(r)

		if err != nil {
			accessDenied(w)
			return
		}

		account, err := s.GetAccountByID(accountId)

		if err != nil {
			accessDenied(w)
			return
		}

		if account.AccountIdentifier != claims["accountIdentifier"] {
			accessDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected string method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func ExtractIdFromRequest(r *http.Request) (int, error) {
	id := mux.Vars(r)["id"]

	parsedId, err := strconv.ParseInt(id, 0, 8)

	if err != nil {
		return 0, err
	}

	return int(parsedId), nil
}
