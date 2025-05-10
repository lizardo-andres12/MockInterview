package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mocker.com/internal/data"
	"go.mocker.com/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	ur *data.UserRepository
}

func NewUserService(ur *data.UserRepository) *UserService {
	return &UserService{ur: ur}
}

func (us *UserService) SignUp(ctx context.Context, user *models.User) error {
	_, err := us.ur.FindByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("Email is already associated with an account, please log in.")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	err = us.ur.Create(ctx, user)
	return err
}

func (us *UserService) LogIn(ctx context.Context, email string, password string) (string, error) {
	user, err := us.ur.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("Email/Password is incorrect.")
	}

	if err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(password),
	); err != nil {
		return "", fmt.Errorf("Email/Password is incorrect.")
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 3),
		},
	)

	tokenString, err := token.SignedString([]byte(os.Getenv("SecretKey")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
