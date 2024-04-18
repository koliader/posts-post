package main

import (
	"context"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
	"github.com/koliader/posts-post.git/internal/util"
	grpc_service "github.com/koliader/posts-post.git/pkg/handler/grpc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	if config.Environment == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	defer connPool.Close()

	store := db.NewStore(connPool)

	go runGrpcServer(config, store)

	server, err := grpc_service.NewServer(config, store)
	if err != nil {
		log.Fatal().Msgf("error creating gRPC server: %v", err)
	}
	err = server.StartConsumer()
	if err != nil {
		log.Fatal().Msgf("error starting RabbitMQ consumer: %v", err)
	}

	select {}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := grpc_service.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}
	// create logger
	grpcLogger := grpc.UnaryInterceptor(grpc_service.GrpcLogger)
	// create server
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterPostServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}
}
