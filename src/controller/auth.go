package controller

import (
	"errors"
	"time"
	"os"

	"go.mocker.com/src/models"
	"go.mocker.com/src/repository"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthController interface {
	Register(email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authController struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewAuthController(repo repository.UserRepository, logger *zap.Logger) AuthController {
	secret := []byte(os.Getenv("JWT_SECRET"))
	return &authController{repo: repo, jwtSecret: secret}
}

func (c *authController) Register(email, password string) (*models.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password required")
	}
	existing, _ := c.repo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already in use")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{Email: email, Password: string(hash)}
	if err := c.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c *authController) Login(email, password string) (string, *models.User, error) {
	user, err := c.repo.GetByEmail(email)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UUID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	})
	signed, err := token.SignedString(c.jwtSecret)
	if err != nil {
		return "", nil, err
	}
	return signed, user, nil
}

