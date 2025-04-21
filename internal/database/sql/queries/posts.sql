-- name: CreatePost :one
INSERT INTO posts (
    user_id, trip_id, title, description, species,
    length_cm, weight_kg, catch_date
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostByIDWithImages :one
SELECT 
    p.*,
    (
        SELECT json_agg(json_build_object(
            'id', pi.id,
            'url', pi.url,
            'is_primary', pi.is_primary,
            'created_at', pi.created_at
        ))
        FROM post_images pi
        WHERE pi.post_id = p.id
    ) AS images
FROM posts p
WHERE p.id = $1;

-- name: ListPostsWithImages :many
SELECT 
    p.*,
    (
        SELECT json_agg(json_build_object(
            'id', pi.id,
            'url', pi.url,
            'is_primary', pi.is_primary,
            'created_at', pi.created_at
        ))
        FROM post_images pi
        WHERE pi.post_id = p.id
    ) AS images
FROM posts p
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListPostsByUserWithImages :many
SELECT 
    p.*,
    (
        SELECT json_agg(json_build_object(
            'id', pi.id,
            'url', pi.url,
            'is_primary', pi.is_primary,
            'created_at', pi.created_at
        ))
        FROM post_images pi
        WHERE pi.post_id = p.id
    ) AS images
FROM posts p
WHERE p.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePost :one
UPDATE posts
SET title = $2,
    description = $3,
    species = $4,
    length_cm = $5,
    weight_kg = $6,
    catch_date = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
