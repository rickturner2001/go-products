package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

type CreateAccountRequest struct {
	Username string `json:"username"`
}

type CreateProductRequestWithId struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`
	Image string `json:"image"`
}

type Account struct {
	ID                int        `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Username          string     `json:"username" gorm:"unique; not null"`
	AccountIdentifier string     `json:"accountIdentifier" gorm:"unique; not null"`
	Products          []*Product `json:"products" gorm:"foreignKey:AccountRefer" `
}

type Product struct {
	ID           int       `json:"id" sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AccountRefer int       `json:"accountRefer"`
	Title        string    `json:"title"`
	Price        float32   `json:"price"`
	Image        string    `json:"image"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewAccount(username string) *Account {
	return &Account{
		Username:          username,
		Products:          []*Product{},
		AccountIdentifier: uuid.New().String(),
	}
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
		ID:        id,
		Title:     title,
		Price:     price,
		Image:     image,
		CreatedAt: time.Now().UTC(),
	}
}
