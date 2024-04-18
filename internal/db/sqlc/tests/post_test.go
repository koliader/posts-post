package db_tests

import (
	"context"
	"testing"

	db "github.com/koliader/posts-post.git/internal/db/sqlc"
	"github.com/koliader/posts-post.git/internal/util"
	"github.com/stretchr/testify/require"
)

func createRandomPost(t *testing.T) db.Post {
	arg := db.CreatePostParams{
		Title: util.RandomString(5),
		Body:  util.RandomString(50),
		Owner: util.RandomEmail(),
	}
	post, err := testStore.CreatePost(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post)

	require.Equal(t, arg.Body, post.Body)
	require.Equal(t, arg.Title, post.Title)
	require.Equal(t, arg.Owner, post.Owner)
	return post
}

func TestCreatePost(t *testing.T) {
	createRandomPost(t)
}

func TestListPosts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomPost(t)
	}
	posts, err := testStore.ListPosts(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, posts)
}

func TestGetPostByTitle(t *testing.T) {
	post1 := createRandomPost(t)
	post2, err := testStore.GetPostByTitle(context.Background(), post1.Title)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, post1.Body, post1.Body)
	require.Equal(t, post1.Title, post1.Title)
	require.Equal(t, post1.Owner, post1.Owner)
}

func TestUpdatePostOwner(t *testing.T) {
	post1 := createRandomPost(t)
	arg := db.UpdatePostOwnerParams{
		Owner: post1.Owner,
		Title: post1.Title,
	}
	post2, err := testStore.UpdatePostOwner(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, post2)

	require.Equal(t, arg.Owner, post2.Owner)
}
