package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secret        []byte
	tokenDuration time.Duration
}

var ErrInvalidToken = errors.New("invalid token")

type UserClaims struct {
	jwt.StandardClaims
	UserID string `json:"sub"`
	Role   string `json:"role"`
}

func NewJWTManager(secret string, tokenDuration time.Duration) (*JWTManager, error) {
	if secret == "" {
		return nil, errors.New("secret is required")
	}
	return &JWTManager{secret: []byte(secret), tokenDuration: tokenDuration}, nil
}

// IssueToken will issue a JWT token with the privided userID as the subject. The token will expire after 15 minutes.
func (s *JWTManager) IssueToken(_ context.Context, user *UserClaims) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.UserID,
		Role:   user.Role,
	}
	// build JWT with necessary claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 基于HMAC签名方法
	return token.SignedString([]byte(s.secret))
}

// ValidateToken will validate the provide JWT against the secret key. It'll then check if the token has expired, and then return the user ID set as the token subject.
func (s *JWTManager) ValidateToken(_ context.Context, accessToken string) (*UserClaims, error) {
	// validate token from the correct secret key and signing method.
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return s.secret, nil
		})
	if err != nil {
		return nil, errors.Join(ErrInvalidToken, err)
	}

	// read claims from payload and extract the user ID.
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
