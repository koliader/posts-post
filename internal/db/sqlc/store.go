package db

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koliader/posts-post.git/internal/util"
)

type Store struct {
	connPool *pgxpool.Pool
	*Queries
	config util.Config
}

func NewStore(connPool *pgxpool.Pool) Store {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	return Store{
		connPool: connPool,
		Queries:  New(connPool),
		config:   config,
	}
}
