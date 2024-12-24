package ports

import (
	"context"

	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
)

type DBPort interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Save(context.Context, *domain.User) error
}
