package domain

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	HashedPassword string `json:"password"`
	Role           string `json:"role"`
	Address        string `json:"address"`
	CreatedAt      int64  `json:"created_at"`
}

func NewUser(name string, email string, phone string, address string, password string, role string) (*User, error) {
	// 系统中不存储明文密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	return &User{
		CreatedAt:      time.Now().Unix(),
		Name:           name,
		HashedPassword: string(hashedPassword),
		Role:           role,
		Email:          email,
		Phone:          phone,
		Address:        address,
	}, nil
}

func (user *User) IsCorrectPassword(password string) bool {
	if user.HashedPassword == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

func (user *User) Clone() *User {
	return &User{
		CreatedAt:      user.CreatedAt,
		Name:           user.Name,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
		Email:          user.Email,
		Phone:          user.Phone,
		Address:        user.Address,
	}
}
