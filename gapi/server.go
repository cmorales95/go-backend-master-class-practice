package gapi

import (
	"fmt"

	db "github.com/cmorales95/go-backend-master-class/db/sqlc"
	"github.com/cmorales95/go-backend-master-class/pb"
	"github.com/cmorales95/go-backend-master-class/token"
	"github.com/cmorales95/go-backend-master-class/util"
)

// Server serves gRPC requests for our banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
