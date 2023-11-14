// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: feedFollow.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE id = $1
`

func (q *Queries) DeleteFeedFollow(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, id)
	return err
}

const follow = `-- name: Follow :one
INSERT INTO feed_follows (
  id, user_id, feed_id, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type FollowParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) Follow(ctx context.Context, arg FollowParams) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, follow,
		arg.ID,
		arg.UserID,
		arg.FeedID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getFeedFollow = `-- name: GetFeedFollow :one
SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetFeedFollow(ctx context.Context, id uuid.UUID) (FeedFollow, error) {
	row := q.db.QueryRowContext(ctx, getFeedFollow, id)
	var i FeedFollow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const getUserFeedFollow = `-- name: GetUserFeedFollow :many
SELECT id, created_at, updated_at, user_id, feed_id FROM feed_follows
WHERE user_id = $1
`

func (q *Queries) GetUserFeedFollow(ctx context.Context, userID uuid.UUID) ([]FeedFollow, error) {
	rows, err := q.db.QueryContext(ctx, getUserFeedFollow, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollow
	for rows.Next() {
		var i FeedFollow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}