-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    inserted.*, 
    users.name AS user_name, 
    feeds.name AS feed_name
FROM inserted
    JOIN users ON users.id = inserted.user_id
    JOIN feeds ON feeds.id = inserted.feed_id;

-- name: GetFeedFollows :many
SELECT 
    feed_follows.*, 
    users.name AS user_name, 
    feeds.name AS feed_name
FROM feed_follows
    JOIN users ON users.id = feed_follows.user_id
    JOIN feeds ON feeds.id = feed_follows.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*, 
    users.name AS user_name, 
    feeds.name AS feed_name
FROM feed_follows
    JOIN users ON users.id = feed_follows.user_id
    JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE users.id = $1;

-- name: DeleteFeedFollow :exec
DELETE 
FROM feed_follows
USING feeds
WHERE feeds.id = feed_follows.feed_id
    AND feed_follows.user_id = $1
    AND feeds.url = $2;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows;
