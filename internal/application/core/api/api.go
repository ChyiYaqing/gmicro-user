package api

import (
	"context"

	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
	"github.com/chyiyaqing/gmicro-user/internal/ports"
)

type Application struct {
	db ports.DBPort // interface as value
}

var _ ports.APIPort = (*Application)(nil)

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a *Application) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	err := a.db.Save(ctx, user)
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

func (a *Application) GetUser(ctx context.Context, user_id int64) (domain.User, error) {
	return a.db.Get(ctx, user_id)
}

func (a *Application) FindUser(ctx context.Context, username string) (domain.User, error) {
	return a.db.Find(ctx, username)
}
