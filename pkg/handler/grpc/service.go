package grpc_service

import (
	"context"

	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
	"google.golang.org/grpc/codes"
)

// * CreatePost

func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostReq) (*pb.CreatePostRes, error) {
	arg := db.CreatePostParams{
		Title: req.GetTitle(),
		Body:  req.GetBody(),
		Owner: req.GetOwnerEmail(),
	}
	post, err := s.store.CreatePost(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, errorResponse(codes.AlreadyExists, "post with this title already created")

		}
		return nil, errorResponse(codes.Unimplemented, "error creating post")
	}
	res := pb.CreatePostRes{
		Post: convertPost(post),
	}
	return &res, nil
}
