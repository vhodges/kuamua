// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Pattern struct {
	ID           pgtype.Int8 `json:"id"`
	PatternName  string      `json:"pattern_name"`
	Pattern      string      `json:"pattern"`
	GroupName    string      `json:"group_name"`
	SubGroupName string      `json:"sub_group_name"`
	OwnerID      string      `json:"owner_id"`
}