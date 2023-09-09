package main

import (
	"time"
)

type CreateProductRequest struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

type CreateProductRequestWithId struct {

	ID string `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

type Product struct {
	ID        int       `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Title     string    `json:"title"`
	Price     float32   `json:"price"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewProduct(title string, price float32, image string) *Product {
	return &Product{
		Title:     title,
		Price:     price,
		Image:     image,
		CreatedAt: time.Now().UTC(),
	}
}

func NewProductWithId(id int, title string, price float32, image string) *Product {
	return &Product{
		ID: 	   id,			
		Title:     title,
		Price:     price,
		Image:     image,
		CreatedAt: time.Now().UTC(),
	}
}

