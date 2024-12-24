package db

import (
	"context"
	"fmt"

	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name    string
	Email   string
	Phone   string
	Address string
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(sqliteDB string) (*Adapter, error) {
	db, openErr := gorm.Open(sqlite.Open(sqliteDB), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db open %v error: %v", sqliteDB, openErr)
	}
	err := db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &Adapter{db: db}, nil
}

func (a *Adapter) Save(ctx context.Context, user *domain.User) error {
	userModel := User{
		Name:    user.Name,
		Email:   user.Email,
		Phone:   user.Phone,
		Address: user.Address,
	}
	res := a.db.Create(&userModel)
	if res.Error == nil {
		user.ID = int64(userModel.ID)
	}
	return res.Error
}

func (a *Adapter) Get(ctx context.Context, id int64) (domain.User, error) {
	var userEntity User
	res := a.db.First(&userEntity, id)
	user := domain.User{
		ID:        int64(userEntity.ID),
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Phone:     userEntity.Phone,
		Address:   userEntity.Address,
		CreatedAt: userEntity.CreatedAt.UnixNano(),
	}
	return user, res.Error
}
