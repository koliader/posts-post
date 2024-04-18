package grpc_service

import (
	"fmt"

	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
	"github.com/koliader/posts-post.git/internal/rabbitmq"
	"github.com/koliader/posts-post.git/internal/util"
)

type Server struct {
	pb.UnimplementedPostServer
	config         util.Config
	store          db.Store
	rabbitmqClient rabbitmq.Client
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	rabbitmqClient, err := rabbitmq.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("error to create rabbitmq client: %v", err)
	}
	err = rabbitmqClient.CreateQueue("updateUserEmail")
	if err != nil {
		return nil, fmt.Errorf("error to create rabbitmq queue: %v", err)
	}

	server := &Server{
		config:         config,
		store:          store,
		rabbitmqClient: *rabbitmqClient,
	}

	return server, nil
}
