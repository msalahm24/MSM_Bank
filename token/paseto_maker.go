package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMaker is PASETO Token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	len := len(symmetricKey)
	if len != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid Key size must be exactly %d char", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(email string,username string, duration time.Duration) (string,*Payload, error) {
	payload, err := NewPayload(email,username, duration)
	if err != nil {
		return "",payload,err
	}
	token,err:=maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token,payload,err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()

	if err != nil {
		return nil, err
	}

	return payload, nil
}
