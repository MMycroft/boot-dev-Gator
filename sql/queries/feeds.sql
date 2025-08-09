-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedByName :one
SELECT 
    feeds.*,
    users.name as user_name
FROM feeds
JOIN users ON users.id = feeds.user_id
WHERE feeds.name = $1
LIMIT 1;

-- name: GetFeedByUrl :one
SELECT 
    feeds.*, 
    users.name as user_name
FROM feeds
JOIN users ON users.id = feeds.user_id
WHERE feeds.url = $1
LIMIT 1;

-- name: GetFeeds :many
SELECT 
    feeds.*, 
    users.name as user_name
FROM feeds
JOIN users ON users.id = feeds.user_id;

-- name: DeleteFeeds :exec
DELETE FROM feeds;
