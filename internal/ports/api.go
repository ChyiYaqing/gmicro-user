package ports

import (
	"context"

	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
)

type APIPort interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
}
