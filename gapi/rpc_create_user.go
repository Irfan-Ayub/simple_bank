package gapi

import (
	"context"

	db "github.com/Irfan-Ayub/simple_bank/db/sqlc"
	"github.com/Irfan-Ayub/simple_bank/pb"
	"github.com/Irfan-Ayub/simple_bank/util"
	"github.com/Irfan-Ayub/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	voilations := validateCreateUserRequest(req)

	if voilations != nil {
		return nil, invalidArgumentError(voilations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Failed to create user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: converUser(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldVoilations("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		voilations = append(voilations, fieldVoilations("password", err))
	}

	if err := val.ValidateFullname(req.GetFullName()); err != nil {
		voilations = append(voilations, fieldVoilations("full_name", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		voilations = append(voilations, fieldVoilations("email", err))
	}

	return voilations
}
