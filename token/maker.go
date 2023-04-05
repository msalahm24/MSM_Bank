package token

import "time"

type Maker interface {
	CreateToken(email string,username string, duration time.Duration) (string,*Payload ,error)
	VerifyToken(token string) (*Payload, error)
}
