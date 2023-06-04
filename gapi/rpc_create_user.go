package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/pb"
	"github.com/mahmoud24598salah/MSM_Bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserRes, error) {
	err := validateReq(req)
	if err != nil {
		return nil, err
	}
	hashedPass, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "can not hash the password")
	}
	arg := db.CreateUserParams{
		Username:   req.GetUsername(),
		HashedPass: hashedPass,
		FullName:   req.GetFullName(),
		Email:      req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user name already exists")
			}
		}

		return nil, status.Errorf(codes.Internal, "can not create user")

	}
	res := &pb.CreateUserRes{
		User: convertUser(user),
	}

	return res, nil
}
