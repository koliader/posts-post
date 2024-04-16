-- name: CreatePost :one
INSERT INTO posts (
  title,
  body,
  owner
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: ListPosts :many
SELECT * FROM posts;

-- name: GetPostByTitle :one
SELECT * FROM posts
WHERE title = $1;