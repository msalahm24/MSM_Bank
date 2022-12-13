package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/util"
	"net/http"
	"time"
)

type createUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}
type createUserResponse struct {
	Username    string    `json:"username"`
	FullName    string    `json:"fullName"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
	PassChanged time.Time `json:"passChanged"`
}

func (server *server) createUser(ctx *gin.Context) {
	var req createUserRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPass, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Username:   req.UserName,
		HashedPass: hashedPass,
		FullName:   req.FullName,
		Email:      req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	res := createUserResponse{
		Username:    user.Username,
		FullName:    user.FullName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		PassChanged: user.PassChanged,
	}
	ctx.JSON(http.StatusOK, res)
}
