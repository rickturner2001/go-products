package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer("0.0.0.0:8080", store)
	server.Run()
}


func MustLoadEnvVariables() (*EnvVariables){
	
	env := &EnvVariables{}

	err := godotenv.Load(); if err != nil {
		log.Fatalf("Could not load env variables: %v", err)
	}

	DB_PASS := 	os.Getenv("DB_PASSWORD")
	ENVIRONMENT := os.Getenv("ENVIRONMENT")

	env.DB_PASSWORD = DB_PASS
	env.ENVIRONMENT = ENVIRONMENT

	return env
}