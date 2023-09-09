package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	router.HandleFunc("/product", makeHTTPHandleFunc(s.handleProduct))
	router.HandleFunc("/product/{id}", makeHTTPHandleFunc(s.handleProduct))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleProduct(w http.ResponseWriter, r *http.Request) error {


	if r.Method == http.MethodGet {
		return s.handleGetProduct(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateProduct(w, r)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteProduct(w, r)
	}

	if r.Method == http.MethodPatch{
		return s.handleUpdateProduct(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetProduct(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["id"]

	if id == "" {
		products, err := s.store.GetAllProducts()

		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, products)
	}


	parsed_id, err := strconv.ParseInt(id, 0, 8)
	
	if err != nil {
		return err
	}
	

	
	product, err := s.store.GetProductByID(int(parsed_id))

	if err != nil {
		return err
	}


	if product == nil {
		return WriteJSON(w, http.StatusNoContent, product)
	}

	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) error {

	product, err := extractProductFromRequest(r)
	
	if err != nil {
		return err
	}

	if err = s.store.CreateProduct(product); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	
	id := mux.Vars(r)["id"]

	parsed_id, err := strconv.ParseInt(id, 0, 8)

	if err != nil {
		return err
	}

	err = s.store.DeleteProduct(int(parsed_id)); if err != nil {
		return err 
	}

	
	return nil
}

func (s *APIServer) handleUpdateProduct(w http.ResponseWriter, r *http.Request) error {
	prod, err := extractProductWithIdFromRequest(r)

	log.Printf("%+v", prod)

	if err != nil{
		return err
	}

	err = s.store.UpdateProduct(prod); if err != nil {
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
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}


func extractProductFromRequest(r *http.Request) (*Product, error){
	createProductReq := new(CreateProductRequest)

	if err := json.NewDecoder(r.Body).Decode(createProductReq); err != nil {
		return nil, err
	}

	parsedPice, err := strconv.Atoi(createProductReq.Price)
	
	if err != nil {
		return nil, err
	}


	product := NewProduct(createProductReq.Title, float32(parsedPice), createProductReq.Image)
	return product, nil
	
}

func extractProductWithIdFromRequest(r *http.Request) (*Product, error){
	createProductReq := new(CreateProductRequestWithId)

	if err := json.NewDecoder(r.Body).Decode(createProductReq); err != nil {
		return nil, err
	}

	parsedPice, err := strconv.Atoi(createProductReq.Price)
	
	if err != nil {
		return nil, err
	}

	parsedId, err := strconv.Atoi(createProductReq.ID)

	product := NewProductWithId(parsedId, createProductReq.Title, float32(parsedPice), createProductReq.Image)

	return product, nil
	

}