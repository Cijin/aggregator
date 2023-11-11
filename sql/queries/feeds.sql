-- name: CreateFeed :one
INSERT INTO feeds (
  id, name, url, user_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListFeed :many
SELECT * FROM feeds
ORDER BY created_at;

-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = $1 LIMIT 1;

-- name: ListFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at
LIMIT $1
OFFSET $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $3
WHERE id = $1;
