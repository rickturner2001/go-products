package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	CreateProduct(*Product) error
	DeleteProduct(int) error
	UpdateProduct(int, *Product) error
	GetProductByID(int) (*Product, error)
	GetAllProducts() ([]*Product, error)

	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *gorm.DB
}


func mustGetDbDnsString() string {
	env := MustLoadEnvVariables()

	log.Printf("Loading .env with environment: %s\n", env.ENVIRONMENT)
	
	if env.ENVIRONMENT == "production" || env.ENVIRONMENT == "" {
		return fmt.Sprintf("user=postgres password=%s host=db.mzxiorcfjuyqnjbrikiz.supabase.co port=5432 dbname=postgres sslmode=disable", env.DB_PASSWORD)
	}else if env.ENVIRONMENT == "development"{
		return "host=product-db user=postgres port=5432 dbname=postgres password=goproduct sslmode=disable"
	}

	log.Fatalf("Invalid environment  [production, development]: %s", env.ENVIRONMENT)
	return ""
	
}

func NewPostgresStore() (*PostgresStore, error) {
	

	dsn := mustGetDbDnsString() 
	

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	err := s.createProductTable()
	if err != nil {
		return err
	}

	err = s.createAccountTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) createProductTable() error {
	return s.db.AutoMigrate(&Product{})
}

func (s *PostgresStore) createAccountTable() error {
	return s.db.AutoMigrate(&Account{})
}

func (s *PostgresStore) CreateProduct(pr *Product) error {
	res := s.db.Create(pr)
	return res.Error

}
func (s *PostgresStore) UpdateProduct(id int, pr *Product) error {

	dbProd, err := s.GetProductByID(id)

	if err != nil {
		return err
	}

	if dbProd == nil {
		return fmt.Errorf("could not find product with id: %v", id)
	}

	res := s.db.Model(&dbProd).Where("id = ?", dbProd.ID).Updates(pr)

	return res.Error
}

func (s *PostgresStore) DeleteProduct(id int) error {
	res := s.db.Delete(&Product{}, id)

	if res.RowsAffected == 0 {
		return fmt.Errorf("could not find product with id: %v", id)
	}

	return res.Error
}
func (s *PostgresStore) GetProductByID(id int) (*Product, error) {
	product := &Product{}

	res := s.db.First(product, id)

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return product, nil
}

func (s *PostgresStore) GetAllProducts() ([]*Product, error) {

	var products []*Product

	res := s.db.Find(&products)

	return products, res.Error
}

func (s *PostgresStore) GetAllAccounts() ([]*Account, error) {
	var accounts []*Account
	res := s.db.Find(&accounts)
	return accounts, res.Error
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	account := &Account{}

	res := s.db.First(account, id)

	if res.RowsAffected == 0 {
		return nil, nil
	}

	return account, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	res := s.db.Create(acc)
	return res.Error
}
