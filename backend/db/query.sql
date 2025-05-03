-- name: CreateFile :one
INSERT INTO files (
  mime_type, file_size, name, ip_addr, hash
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: GetFileFromHash :one
SELECT id, mime_type FROM files
WHERE hash = $1;

-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1;
