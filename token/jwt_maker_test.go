package token

import ( 
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mahmoud24598salah/MSM_Bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)
	userName := util.RandomOwnerName()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userName,duration)
	require.NoError(t, err)
	require.NotEmpty(t,token)
	payload, err := maker.VerifyToken(token)
	require.NoError(t,err)
	require.NotEmpty(t,payload)
	require.NotZero(t,payload.ID)
	require.Equal(t,userName,payload.UserName)
	require.WithinDuration(t,issuedAt,payload.IssuedAt,time.Second)
	require.WithinDuration(t,expiredAt,payload.ExpiredAt,time.Second)
}


func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t,err)
	token, err := maker.CreateToken(util.RandomOwnerName(),-time.Minute)
	require.NoError(t,err)
	payload, err := maker.VerifyToken(token)
	require.Error(t,err )
	require.EqualError(t,err,ErrExpireToken.Error())
	require.Nil(t,payload)
}


func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload,err := NewPayload(util.RandomOwnerName(),time.Minute)
	require.NoError(t, err)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone,payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t,err)
	maker, err :=NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	payload,err = maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
