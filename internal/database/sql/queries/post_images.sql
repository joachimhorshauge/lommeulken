-- name: AddPostImage :one
INSERT INTO post_images (post_id, url, is_primary)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetImagesByPostID :many
SELECT * FROM post_images
WHERE post_id = $1
ORDER BY created_at ASC;

-- name: GetPrimaryImageForPost :one
SELECT * FROM post_images
WHERE post_id = $1 AND is_primary = TRUE
LIMIT 1;

-- name: DeleteImageByID :exec
DELETE FROM post_images
WHERE id = $1;

-- name: DeleteImagesByPostID :exec
DELETE FROM post_images
WHERE post_id = $1;

