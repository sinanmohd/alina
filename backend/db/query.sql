-- name: CreateFile :exec
INSERT INTO files (
  mime_type, file_size, name
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files WHERE id = ?;
