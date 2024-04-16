package grpc_service

import (
	"context"

	"github.com/jackc/pgx/v5"
	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
	"google.golang.org/grpc/codes"
)

// * CreatePost

func (s *Server) CreatePost(ctx context.Context, req *pb.CreatePostReq) (*pb.CreatePostRes, error) {
	arg := db.CreatePostParams{
		Title: req.GetTitle(),
		Body:  req.GetBody(),
		Owner: req.GetOwner(),
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

// * GetPost
func (s *Server) GetPost(ctx context.Context, req *pb.GetPostReq) (*pb.GetPostRes, error) {
	post, err := s.store.GetPostByTitle(ctx, req.Title)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errorResponse(codes.NotFound, "post not found")
		}
		return nil, errorResponse(codes.Unimplemented, "error to get post")
	}
	res := pb.GetPostRes{
		Post: convertPost(post),
	}
	return &res, nil
}

// * ListPosts
func (s *Server) ListPosts(ctx context.Context, req *pb.Empty) (*pb.ListPostsRes, error) {
	var converted []*pb.PostEntity

	posts, err := s.store.ListPosts(ctx)
	if err != nil {
		return nil, errorResponse(codes.Unimplemented, "error to list posts")
	}
	for _, post := range posts {
		converted = append(converted, convertPost(post))
	}
	res := pb.ListPostsRes{
		Posts: converted,
	}
	return &res, nil
}
