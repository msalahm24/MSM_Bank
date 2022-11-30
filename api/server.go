package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
)

 type server struct{
	store db.Store
	router *gin.Engine 
 }

 func Newserver(store db.Store) *server{
	server := &server{store: store}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok{
		v.RegisterValidation("currency",validCurrency)
	}
	router.POST("/accounts",server.createAccount)
	router.GET("/accounts/:id",server.getAccount)
	router.GET("/accounts",server.listAccounts) 
	router.DELETE("/accounts",server.deleteAccount)
	router.PUT("/accounts",server.putAccount)
	router.POST("/transfer",server.CreateTransfer)
	server.router = router

	return server
 }

func(server *server) Start(address string)error{
	return server.router.Run(address)
}

 func errResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}