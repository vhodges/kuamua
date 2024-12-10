-- name: GetPattern :one
SELECT * FROM kuamua_patterns
WHERE id = $1 LIMIT 1;

-- name: ListOwnerGroupPatterns :many
SELECT * FROM kuamua_patterns
WHERE owner_id = $1
  AND group_name = $2
  AND sub_group_name = $3
ORDER BY pattern_name;

-- name: CreatePattern :one
INSERT INTO kuamua_patterns (
  pattern_name, pattern, group_name, sub_group_name, owner_id
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdatePattern :exec
UPDATE kuamua_patterns
  set pattern_name = $2,
  pattern = $3,
  group_name = $4,
  sub_group_name = $5,
  owner_id = $6
WHERE id = $1;

-- name: DeletePattern :exec
DELETE FROM kuamua_patterns
WHERE id = $1;