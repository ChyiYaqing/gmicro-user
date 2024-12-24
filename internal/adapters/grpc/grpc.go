package grpc

import (
	"context"

	"github.com/chyiyaqing/gmicro-proto/golang/user"
	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
	log "github.com/sirupsen/logrus"
)

func (a *Adapter) Create(ctx context.Context, request *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	log.WithContext(ctx).Info("Creating user...")
	newUser := domain.NewUser(request.Name, request.Email, request.Phone, request.Address)
	result, err := a.api.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &user.CreateUserResponse{
		UserId: result.ID,
	}, nil
}

func (a *Adapter) Get(ctx context.Context, request *user.GetUserRequest) (*user.GetUserResponse, error) {
	log.WithContext(ctx).Info("Get user...")
	result, err := a.api.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &user.GetUserResponse{
		UserId:  result.ID,
		Name:    result.Name,
		Email:   result.Email,
		Phone:   result.Phone,
		Address: result.Address,
	}, nil
}
