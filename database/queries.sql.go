// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPattern = `-- name: CreatePattern :one
INSERT INTO kuamua_patterns (
  pattern_name, pattern, group_name, sub_group_name, owner_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, pattern_name, pattern, group_name, sub_group_name, owner_id
`

type CreatePatternParams struct {
	PatternName  string `json:"pattern_name"`
	Pattern      string `json:"pattern"`
	GroupName    string `json:"group_name"`
	SubGroupName string `json:"sub_group_name"`
	OwnerID      string `json:"owner_id"`
}

func (q *Queries) CreatePattern(ctx context.Context, arg CreatePatternParams) (Pattern, error) {
	row := q.db.QueryRow(ctx, createPattern,
		arg.PatternName,
		arg.Pattern,
		arg.GroupName,
		arg.SubGroupName,
		arg.OwnerID,
	)
	var i Pattern
	err := row.Scan(
		&i.ID,
		&i.PatternName,
		&i.Pattern,
		&i.GroupName,
		&i.SubGroupName,
		&i.OwnerID,
	)
	return i, err
}

const deletePattern = `-- name: DeletePattern :exec
DELETE FROM kuamua_patterns
WHERE id = $1
`

func (q *Queries) DeletePattern(ctx context.Context, id pgtype.Int8) error {
	_, err := q.db.Exec(ctx, deletePattern, id)
	return err
}

const getPattern = `-- name: GetPattern :one
SELECT id, pattern_name, pattern, group_name, sub_group_name, owner_id FROM kuamua_patterns
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPattern(ctx context.Context, id pgtype.Int8) (Pattern, error) {
	row := q.db.QueryRow(ctx, getPattern, id)
	var i Pattern
	err := row.Scan(
		&i.ID,
		&i.PatternName,
		&i.Pattern,
		&i.GroupName,
		&i.SubGroupName,
		&i.OwnerID,
	)
	return i, err
}

const listOwnerGroupPatterns = `-- name: ListOwnerGroupPatterns :many
SELECT id, pattern_name, pattern, group_name, sub_group_name, owner_id FROM kuamua_patterns
WHERE owner_id = $1
  AND group_name = $2
  AND sub_group_name = $3
ORDER BY pattern_name
`

type ListOwnerGroupPatternsParams struct {
	OwnerID      string `json:"owner_id"`
	GroupName    string `json:"group_name"`
	SubGroupName string `json:"sub_group_name"`
}

func (q *Queries) ListOwnerGroupPatterns(ctx context.Context, arg ListOwnerGroupPatternsParams) ([]Pattern, error) {
	rows, err := q.db.Query(ctx, listOwnerGroupPatterns, arg.OwnerID, arg.GroupName, arg.SubGroupName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pattern
	for rows.Next() {
		var i Pattern
		if err := rows.Scan(
			&i.ID,
			&i.PatternName,
			&i.Pattern,
			&i.GroupName,
			&i.SubGroupName,
			&i.OwnerID,
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

const updatePattern = `-- name: UpdatePattern :exec
UPDATE kuamua_patterns
  set pattern_name = $2,
  pattern = $3,
  group_name = $4,
  sub_group_name = $5,
  owner_id = $6
WHERE id = $1
`

type UpdatePatternParams struct {
	ID           pgtype.Int8 `json:"id"`
	PatternName  string      `json:"pattern_name"`
	Pattern      string      `json:"pattern"`
	GroupName    string      `json:"group_name"`
	SubGroupName string      `json:"sub_group_name"`
	OwnerID      string      `json:"owner_id"`
}

func (q *Queries) UpdatePattern(ctx context.Context, arg UpdatePatternParams) error {
	_, err := q.db.Exec(ctx, updatePattern,
		arg.ID,
		arg.PatternName,
		arg.Pattern,
		arg.GroupName,
		arg.SubGroupName,
		arg.OwnerID,
	)
	return err
}
