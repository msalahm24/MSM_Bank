package gapi

import (
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User)*pb.User{
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PassChangedAt: timestamppb.New(user.PassChanged),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}