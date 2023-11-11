-- +goose Up
CREATE TABLE IF NOT EXISTS feed_follows (
  id UUID NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL,
  feed_id UUID NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (feed_id) REFERENCES feeds(id),
  UNIQUE (feed_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS feed_follows;
