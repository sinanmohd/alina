-- name: FileCreate :one
INSERT INTO files (
  mime_type, file_size, name, ip_addr, hash
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: ChunkedCreate :one
INSERT INTO chunked (
  file_size, name, ip_addr, chunks_left, chunks_total
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id;

-- name: ChunkedLeftDecrement :one
UPDATE chunked
SET chunks_left = GREATEST(0, chunks_left - 1), last_access = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING chunks_left;

-- name: ChunkedDelete :exec
DELETE FROM chunked
WHERE id = $1;

-- name: ChunkedFromId :one
SELECT file_size, chunks_total, created_at, last_access FROM chunked
WHERE id = $1;

-- name: ChunkedToFile :one
INSERT INTO files (mime_type, file_size, name, ip_addr, hash)
SELECT $2, chunked.file_size, chunked.name, chunked.ip_addr, $3 
FROM chunked
WHERE chunked.id = $1
RETURNING files.id;

-- name: FileFromHash :one
SELECT id, mime_type FROM files
WHERE hash = $1;

-- name: FileDelete :exec
DELETE FROM files
WHERE id = $1;
