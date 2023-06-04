package gapi

import (
	"context"
	"database/sql"

	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/pb"
	"github.com/mahmoud24598salah/MSM_Bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserRes, error) {

	user, err := server.store.GetUser(ctx, req.GetEmail())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "this user not exsits")
		}
		return nil, status.Errorf(codes.Internal, "can not get the user")
	}
	err = util.CheckPass(user.HashedPass, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "incorrect pass")
	}
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Email, user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not create token")
	}
	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Email, user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not create redresh token")
	}
	userAgent, clientIP := getMetaData(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIP,
		IsBlocked:    false,
		ExpiresIt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not create session")
	}
	res := pb.LoginUserRes{
		User:                  convertUser(user),
		Session_ID:            session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
	}

	return &res, nil
}
