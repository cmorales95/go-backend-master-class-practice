package gapi

import (
	"context"
	"database/sql"

	db "github.com/cmorales95/go-backend-master-class/db/sqlc"
	"github.com/cmorales95/go-backend-master-class/pb"
	"github.com/cmorales95/go-backend-master-class/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed validating user: %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "incorrect password")
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(user.Username, s.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token:")
	}

	metadata := s.extractMetadata(ctx)
	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed creating session: %s", err)
	}

	return &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}, nil
}
