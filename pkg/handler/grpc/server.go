package grpc_service

import (
	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
	"github.com/koliader/posts-post.git/internal/util"
)

type Server struct {
	pb.UnimplementedPostServer
	config util.Config
	store  db.Store
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	return server, nil
}
