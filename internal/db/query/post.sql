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

-- name: ListPostsByOwner :many
SELECT * FROM posts
WHERE owner = $1;

-- name: GetPostByTitle :one
SELECT * FROM posts
WHERE title = $1;

-- name: UpdatePostOwner :one
UPDATE posts
SET owner = $2
WHERE title = $1
RETURNING *;