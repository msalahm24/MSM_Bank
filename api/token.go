package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expire_at"`
}

func (server *server) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
	}

	if session.IsBlocked {
		err = errors.New("blocked session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	if session.Username != refreshPayload.UserName {
		err = errors.New("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	if session.RefreshToken != req.RefreshToken {
		err = errors.New("mismatched session token")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	if time.Now().After(session.ExpiresIt) {
		err = errors.New("expired session")
		ctx.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshPayload.Email, refreshPayload.UserName, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	res := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessTokenPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, res)
}
