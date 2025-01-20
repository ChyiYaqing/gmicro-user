package grpc

import (
	"context"

	userv1 "github.com/chyiyaqing/gmicro-proto/golang/user/v1"
	"github.com/chyiyaqing/gmicro-user/internal/adapters/auth"
	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Adapter) Login(ctx context.Context, request *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := a.api.FindUser(ctx, request.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}

	if user.ID == 0 || !user.IsCorrectPassword(request.Password) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}
	token, err := a.jwtManager.IssueToken(ctx, &auth.UserClaims{UserID: string(user.ID), Role: user.Role})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}
	res := &userv1.LoginResponse{AccessToken: token}
	return res, nil
}

func (a *Adapter) Create(ctx context.Context, request *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	log.WithContext(ctx).Info("Creating user...")
	newUser, err := domain.NewUser(request.Username, request.Email, request.Phone, request.Address, request.Password, request.Role)
	result, err := a.api.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{
		UserId: result.ID,
	}, nil
}

func (a *Adapter) Get(ctx context.Context, request *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	log.WithContext(ctx).Info("Get user...")
	// extract user ID from context
	userID, ok := ctx.Value("user_id").(string)
	if ok {
		log.Infof("get userID: %s", userID)
	}
	// reset of your RPC implementation here...
	result, err := a.api.GetUser(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &userv1.GetUserResponse{
		UserId:  result.ID,
		Name:    result.Name,
		Email:   result.Email,
		Phone:   result.Phone,
		Address: result.Address,
	}, nil
}
