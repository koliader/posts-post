package grpc_service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/pb"
	"github.com/koliader/posts-post.git/internal/rabbitmq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
)

// * CreatePost

// TODO check user is exists before to create post
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

func (s *Server) ListPostsByUser(ctx context.Context, req *pb.ListPostsByUserReq) (*pb.ListPostsRes, error) {
	var converted []*pb.PostEntity

	posts, err := s.store.ListPostsByOwner(ctx, req.GetOwner())
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

func (s *Server) StartConsumer() error {
	msgs, err := s.rabbitmqClient.Channel.Consume(
		"updateUserEmail",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error to consume messages")
	}
	go func() {
		for d := range msgs {
			var msgBody rabbitmq.UpdateEmailMessage
			err := json.Unmarshal(d.Body, &msgBody)
			if err != nil {
				log.Fatal().Msg(fmt.Sprintf("Error unmarshalling message: %v", err))
				continue
			}
			userPosts, err := s.store.ListPostsByOwner(context.Background(), msgBody.Email)
			if err != nil {
				log.Fatal().Msg("error to list user posts")
			}
			for _, post := range userPosts {
				arg := db.UpdatePostOwnerParams{
					Owner: msgBody.NewEmail,
					Title: post.Title,
				}
				_, err = s.store.UpdatePostOwner(context.Background(), arg)
				if err != nil {
					log.Fatal().Msg(fmt.Sprintf("Error to update post owner: %v", err))
				}
			}
		}
	}()
	return nil
}
