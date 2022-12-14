package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTMaker struct {
	Secret string
}

const MIN_SECRET_KEY_SIZE = 32

func NewJWTMaker(secretkey string) (IMaker, error) {
	if len(secretkey) < MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalied secret key size must be at least %d char", MIN_SECRET_KEY_SIZE)
	}
	return &JWTMaker{secretkey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

}
