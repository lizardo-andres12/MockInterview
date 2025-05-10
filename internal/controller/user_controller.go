package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mocker.com/internal/models"
	"go.mocker.com/internal/service"
)

type UserController struct {
	us *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{us: us}
}

func (uc *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var user models.User
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if dec.More() {
		http.Error(w, "JSON contained unexpected fields", http.StatusUnsupportedMediaType)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second * 3)
	defer cancel()

	err := uc.us.SignUp(ctx, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"message":"user created"}`))
}
