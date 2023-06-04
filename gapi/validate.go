package gapi

import (
	"github.com/asaskevich/govalidator"
	"github.com/mahmoud24598salah/MSM_Bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func validateReq(req *pb.CreateUserRequest) error {
	if req.GetUsername() == "" {
		return status.Errorf(codes.InvalidArgument, "invalid user name")
	}
	if req.GetFullName() == "" {
		return status.Errorf(codes.InvalidArgument, "invalid full name")
	}
	if !govalidator.IsEmail(req.GetEmail()) {
		return status.Errorf(codes.InvalidArgument, "invalid email")
	}
	if req.GetPassword() == "" || len(req.GetPassword()) < 8{
		return status.Errorf(codes.InvalidArgument, "invalid password")
	}
	return nil
}
