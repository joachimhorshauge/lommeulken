// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: posts.sql

package dbstore

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
    user_id, trip_id, title, description, species,
    length_cm, weight_kg, catch_date
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, trip_id, title, description, species, length_cm, weight_kg, catch_date, created_at, updated_at
`

type CreatePostParams struct {
	UserID      uuid.UUID
	TripID      pgtype.UUID
	Title       string
	Description pgtype.Text
	Species     string
	LengthCm    pgtype.Int4
	WeightKg    pgtype.Float8
	CatchDate   pgtype.Timestamptz
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, createPost,
		arg.UserID,
		arg.TripID,
		arg.Title,
		arg.Description,
		arg.Species,
		arg.LengthCm,
		arg.WeightKg,
		arg.CatchDate,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TripID,
		&i.Title,
		&i.Description,
		&i.Species,
		&i.LengthCm,
		&i.WeightKg,
		&i.CatchDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deletePost, id)
	return err
}

const getPostByIDWithImages = `-- name: GetPostByIDWithImages :one
SELECT 
    p.id, p.user_id, p.trip_id, p.title, p.description, p.species, p.length_cm, p.weight_kg, p.catch_date, p.created_at, p.updated_at,
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
WHERE p.id = $1
`

type GetPostByIDWithImagesRow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	TripID      pgtype.UUID
	Title       string
	Description pgtype.Text
	Species     string
	LengthCm    pgtype.Int4
	WeightKg    pgtype.Float8
	CatchDate   pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	Images      []byte
}

func (q *Queries) GetPostByIDWithImages(ctx context.Context, id uuid.UUID) (GetPostByIDWithImagesRow, error) {
	row := q.db.QueryRow(ctx, getPostByIDWithImages, id)
	var i GetPostByIDWithImagesRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TripID,
		&i.Title,
		&i.Description,
		&i.Species,
		&i.LengthCm,
		&i.WeightKg,
		&i.CatchDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Images,
	)
	return i, err
}

const listPostsByUserWithImages = `-- name: ListPostsByUserWithImages :many
SELECT 
    p.id, p.user_id, p.trip_id, p.title, p.description, p.species, p.length_cm, p.weight_kg, p.catch_date, p.created_at, p.updated_at,
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
LIMIT $2 OFFSET $3
`

type ListPostsByUserWithImagesParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

type ListPostsByUserWithImagesRow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	TripID      pgtype.UUID
	Title       string
	Description pgtype.Text
	Species     string
	LengthCm    pgtype.Int4
	WeightKg    pgtype.Float8
	CatchDate   pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	Images      []byte
}

func (q *Queries) ListPostsByUserWithImages(ctx context.Context, arg ListPostsByUserWithImagesParams) ([]ListPostsByUserWithImagesRow, error) {
	rows, err := q.db.Query(ctx, listPostsByUserWithImages, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPostsByUserWithImagesRow
	for rows.Next() {
		var i ListPostsByUserWithImagesRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TripID,
			&i.Title,
			&i.Description,
			&i.Species,
			&i.LengthCm,
			&i.WeightKg,
			&i.CatchDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Images,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listPostsWithImages = `-- name: ListPostsWithImages :many
SELECT 
    p.id, p.user_id, p.trip_id, p.title, p.description, p.species, p.length_cm, p.weight_kg, p.catch_date, p.created_at, p.updated_at,
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
WHERE 
    ($1::boolean = FALSE OR p.user_id = $2)
    AND (
        $3::boolean = FALSE OR 
        p.species = ANY($4::VARCHAR[])
    )
ORDER BY
    CASE WHEN $5::text = 'created_at' AND $6::text = 'asc' THEN p.created_at ELSE NULL END ASC NULLS LAST,
    CASE WHEN $5::text = 'created_at' AND $6::text = 'desc' THEN p.created_at ELSE NULL END DESC NULLS LAST,
    CASE WHEN $5::text = 'length_cm' AND $6::text = 'asc' THEN p.length_cm ELSE NULL END ASC NULLS LAST,
    CASE WHEN $5::text = 'length_cm' AND $6::text = 'desc' THEN p.length_cm ELSE NULL END DESC NULLS LAST,
    CASE WHEN $5::text = 'weight_kg' AND $6::text = 'asc' THEN p.weight_kg ELSE NULL END ASC NULLS LAST,
    CASE WHEN $5::text = 'weight_kg' AND $6::text = 'desc' THEN p.weight_kg ELSE NULL END DESC NULLS LAST,
    p.created_at DESC 
LIMIT $8 OFFSET $7
`

type ListPostsWithImagesParams struct {
	FilterByUserID  bool
	UserID          uuid.UUID
	FilterBySpecies bool
	Species         []string
	SortColumn      string
	SortDirection   string
	ResultOffset    int32
	ResultLimit     int32
}

type ListPostsWithImagesRow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	TripID      pgtype.UUID
	Title       string
	Description pgtype.Text
	Species     string
	LengthCm    pgtype.Int4
	WeightKg    pgtype.Float8
	CatchDate   pgtype.Timestamptz
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	Images      []byte
}

func (q *Queries) ListPostsWithImages(ctx context.Context, arg ListPostsWithImagesParams) ([]ListPostsWithImagesRow, error) {
	rows, err := q.db.Query(ctx, listPostsWithImages,
		arg.FilterByUserID,
		arg.UserID,
		arg.FilterBySpecies,
		arg.Species,
		arg.SortColumn,
		arg.SortDirection,
		arg.ResultOffset,
		arg.ResultLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListPostsWithImagesRow
	for rows.Next() {
		var i ListPostsWithImagesRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TripID,
			&i.Title,
			&i.Description,
			&i.Species,
			&i.LengthCm,
			&i.WeightKg,
			&i.CatchDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Images,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET title = $2,
    description = $3,
    species = $4,
    length_cm = $5,
    weight_kg = $6,
    catch_date = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, trip_id, title, description, species, length_cm, weight_kg, catch_date, created_at, updated_at
`

type UpdatePostParams struct {
	ID          uuid.UUID
	Title       string
	Description pgtype.Text
	Species     string
	LengthCm    pgtype.Int4
	WeightKg    pgtype.Float8
	CatchDate   pgtype.Timestamptz
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRow(ctx, updatePost,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Species,
		arg.LengthCm,
		arg.WeightKg,
		arg.CatchDate,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TripID,
		&i.Title,
		&i.Description,
		&i.Species,
		&i.LengthCm,
		&i.WeightKg,
		&i.CatchDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
