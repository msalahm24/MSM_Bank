package gapi

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/pb"
	"github.com/mahmoud24598salah/MSM_Bank/token"
	"github.com/mahmoud24598salah/MSM_Bank/util"
)

type server struct {
	pb.UnimplementedMSMBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func Newserver(config util.Config, store db.Store) (*server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricLKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}


func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
