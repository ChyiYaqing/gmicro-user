package db

import (
	"context"
	"fmt"

	"github.com/chyiyaqing/gmicro-user/internal/application/core/domain"
	"github.com/chyiyaqing/gmicro-user/internal/ports"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string `gorm:"unique;not null"`
	Email          string `gorm:"unique;not null"`
	Phone          string
	HashedPassword string `gorm:"size:60;not null"` // bcrypt 哈希通常是 60 字节
	Role           string
	Address        string
}

type Adapter struct {
	db *gorm.DB
}

var _ ports.DBPort = (*Adapter)(nil)

func NewAdapter(sqliteDB string) (*Adapter, error) {
	db, openErr := gorm.Open(sqlite.Open(sqliteDB), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db open %v error: %v", sqliteDB, openErr)
	}
	err := db.AutoMigrate(&User{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	// 检查是否已存在管理员账户
	var count int64
	db.Model(&User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to generate hash password: %v", err)
		}
		adminUser := User{
			Name:           "admin",
			Email:          "admin@example.com",
			Phone:          "",
			HashedPassword: string(hashedPassword),
			Role:           "admin",
			Address:        "",
		}
		result := db.Create(&adminUser)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to create admin user: %v", result.Error)
		}
	}

	return &Adapter{db: db}, nil
}

func (a *Adapter) Save(ctx context.Context, user *domain.User) error {
	userModel := User{
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
		Address:        user.Address,
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
		Role:      userEntity.Role,
		Address:   userEntity.Address,
		CreatedAt: userEntity.CreatedAt.UnixNano(),
	}
	return user, res.Error
}

func (a *Adapter) Find(ctx context.Context, username string) (domain.User, error) {
	var userEntity User
	res := a.db.First(&userEntity, "name = ?", username)
	if res.Error != nil {
		return domain.User{}, res.Error
	}

	user := domain.User{
		ID:             int64(userEntity.ID),
		Name:           userEntity.Name,
		Email:          userEntity.Email,
		Phone:          userEntity.Phone,
		HashedPassword: userEntity.HashedPassword,
		Role:           userEntity.Role,
		Address:        userEntity.Address,
		CreatedAt:      userEntity.CreatedAt.UnixNano(),
	}
	return user, res.Error
}
