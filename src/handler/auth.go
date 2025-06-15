package handler

import (
	"context"
	"go.mocker.com/src/controller"
	"go.mocker.com/src/proto"
)

type AuthHandler struct {
	ctrl controller.AuthController
	proto.UnimplementedAuthServiceServer
}

func NewAuthHandler(ctrl controller.AuthController) *AuthHandler {
	return &AuthHandler{ctrl: ctrl}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	req *proto.RegisterRequest,
) (*proto.UserResponse, error) {
	user, err := h.ctrl.Register(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.UserResponse{Uuid: user.UUID, Email: user.Email}, nil
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *proto.LoginRequest,
) (*proto.LoginResponse, error) {
	token, user, err := h.ctrl.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.LoginResponse{
		Token: token,
		User: &proto.UserResponse{
			Uuid: user.UUID,
			Email: user.Email,
		},
	}, nil
}

