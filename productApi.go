package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) HandleProduct(w http.ResponseWriter, r *http.Request) error {

	if r.Method == http.MethodGet {
		return s.handleGetProduct(w, r)
	}

	if r.Method == http.MethodPost {
		return s.handleCreateProduct(w, r)
	}

	if r.Method == http.MethodDelete {
		return s.handleDeleteProduct(w, r)
	}

	if r.Method == http.MethodPatch {
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

	parsedId, err := strconv.ParseInt(id, 0, 8)

	if err != nil {
		return err
	}

	product, err := s.store.GetProductByID(int(parsedId))

	if err != nil {
		return err
	}

	if product == nil {
		return WriteJSON(w, http.StatusNotFound, ApiError{Error: fmt.Sprintf("Could not find product with id: %d", parsedId)})
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

	parsedId, err := ExtractIdFromRequest(r)

	err = s.store.DeleteProduct(int(parsedId))
	if err != nil {
		return WriteJSON(w, http.StatusNotFound, ApiError{Error: err.Error()})
	}

	return nil
}

func (s *APIServer) handleUpdateProduct(w http.ResponseWriter, r *http.Request) error {

	id, err := ExtractIdFromRequest(r)

	if err != nil {
		return err
	}

	prod, err := extractProductFromRequest(r)

	if err != nil {
		return err
	}

	err = s.store.UpdateProduct(id, prod)
	if err != nil {
		return err
		// return WriteJSON(w, )
	}

	prod.ID = id
	return WriteJSON(w, http.StatusOK, prod)
}

func extractProductFromRequest(r *http.Request) (*Product, error) {
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

func extractProductWithIdFromRequest(r *http.Request) (*Product, error) {
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
