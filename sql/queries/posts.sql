-- name: CreatePost :one
INSERT INTO posts (
  id, created_at, updated_at, title, url, description, published_at, feed_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListPostsByUser :many
SELECT p.* FROM posts p
INNER JOIN feed_follows f 
ON p.feed_id = f.feed_id
AND f.user_id = $1
ORDER BY p.published_at DESC;
