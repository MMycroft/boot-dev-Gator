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

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = NOW(),
    last_fetched_at = NOW()
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY feeds.updated_at ASC NULLS FIRST
LIMIT 1;

-- name: DeleteFeeds :exec
DELETE FROM feeds;
