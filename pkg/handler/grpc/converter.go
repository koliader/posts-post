package grpc_service

import (
	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
)

func convertPost(post db.Post) *pb.PostEntity {
	converted := pb.PostEntity{
		Title:      post.Title,
		Body:       post.Body,
		OwnerEmail: post.Owner,
	}
	return &converted
}
