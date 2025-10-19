package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/Irfan-Ayub/simple_bank/db/sqlc"
	"github.com/Irfan-Ayub/simple_bank/pb"
	"github.com/Irfan-Ayub/simple_bank/util"
	"github.com/Irfan-Ayub/simple_bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	voilations := validateUpdateUserRequest(req)

	if voilations != nil {
		return nil, invalidArgumentError(voilations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's account")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to udpate user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: converUser(user),
	}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (voilations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		voilations = append(voilations, fieldVoilations("username", err))
	}

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			voilations = append(voilations, fieldVoilations("password", err))
		}
	}

	if req.FullName != nil {
		if err := val.ValidateFullname(req.GetFullName()); err != nil {
			voilations = append(voilations, fieldVoilations("full_name", err))
		}
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			voilations = append(voilations, fieldVoilations("email", err))
		}
	}

	return voilations
}
