package domain

import "time"

type User struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	CreatedAt int64  `json:"created_at"`
}

func NewUser(name string, email string, phone string, address string) User {
	return User{
		CreatedAt: time.Now().Unix(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
	}
}
