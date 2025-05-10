package main

import (
	"fmt"
	"log"
	"net/http"

	"go.mocker.com/internal/controller"
	"go.mocker.com/internal/data"
	"go.mocker.com/internal/service"
)

func main() {
	config, err := data.GenerateConfig()
	if err != nil {
		log.Fatalf("%v: failed", err)
	}

	db, err := data.NewDatabase(config)
	if err != nil {
		log.Fatalf("%v: failed", err)
	}

	ur := data.NewUserRepository(db)
	us := service.NewUserService(ur)
	uc := controller.NewUserController(us)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", uc.SignUp)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "welcome to mocker!")
	})

	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("%v: failed", err)
	}
}

