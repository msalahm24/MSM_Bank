package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/token"
	"github.com/mahmoud24598salah/MSM_Bank/util"
)

type server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.IMaker
	router     *gin.Engine
}

func Newserver(config util.Config, store db.Store) (*server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricLKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setUpRouter()
	return server, nil
}

func (server *server) setUpRouter() {
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.POST("/users", server.createUser)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.DELETE("/accounts", server.deleteAccount)
	router.PUT("/accounts", server.putAccount)
	router.POST("/transfer", server.CreateTransfer)
	router.POST("/users/login", server.loginUser)
	server.router = router
}

func (server *server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
