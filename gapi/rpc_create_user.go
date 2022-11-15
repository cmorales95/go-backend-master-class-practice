package gapi

import (
	"context"

	db "github.com/cmorales95/go-backend-master-class/db/sqlc"
	"github.com/cmorales95/go-backend-master-class/pb"
	"github.com/cmorales95/go-backend-master-class/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to has password %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: ", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "failed to create user: ", err)
	}

	return &pb.CreateUserResponse{
		User: convertUser(user),
	}, nil
}
