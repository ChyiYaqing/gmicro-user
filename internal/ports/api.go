package ports

import (
	"context"

	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
)

type APIPort interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	FindUser(ctx context.Context, username string) (domain.User, error)
}
