package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretkey string
}

const MIN_SECRET_KEY_SIZE = 32

func NewJWTMaker(secretkey string) (Maker, error) {
	if len(secretkey) < MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalied secret key size must be at least %d char", MIN_SECRET_KEY_SIZE)
	}
	return &JWTMaker{secretkey}, nil
}


func (maker *JWTMaker) CreateToken(email string,username string, duration time.Duration) (string, *Payload,error) {
	payload, err := NewPayload(email,username, duration)
	if err != nil {
		return "", payload,err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token,err:=jwtToken.SignedString([]byte(maker.secretkey))
	return token,payload,err
}


func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error){
		_,ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok{
			return nil,ErrInvalidToken
		}
		return []byte(maker.secretkey),nil
	}
	jwtToken, err := jwt.ParseWithClaims(token,&Payload{},keyFunc)
	if err != nil{
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpireToken){
			return nil, ErrExpireToken
		}
		return nil,ErrInvalidToken
	}
	
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok{
		return nil, ErrInvalidToken
	}
	return payload, nil
}
