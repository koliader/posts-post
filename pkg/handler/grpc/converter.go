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

func convertPosts(posts []db.Post) []*pb.PostEntity {
	var converted []*pb.PostEntity
	for _, post := range posts {
		converted = append(converted, convertPost(post))
	}
	return converted
}
