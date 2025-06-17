-- name: UserAgentIdGet :one
WITH res AS (
    INSERT INTO user_agents (user_agent)
    VALUES ($1)
    ON CONFLICT (user_agent) DO NOTHING
    RETURNING id
)
SELECT id
FROM res
UNION
SELECT id
FROM user_agents
WHERE user_agent = $1
LIMIT 1;

INSERT INTO user_agents (user_agent)
SELECT $1
WHERE NOT EXISTS (
  SELECT id
  FROM user_agents
  WHERE user_agent = $1
)
RETURNING id;

-- name: UploadCreate :one
INSERT INTO uploads (
  ip_addr, user_agent, file, name
) VALUES (
  $1, $2, $3, $4
)
RETURNING id;

-- name: FileCreate :one
INSERT INTO files (
  mime_type, file_size, hash
) VALUES (
  $1, $2, $3
)
RETURNING id;

-- name: ChunkedCreate :one
INSERT INTO chunked (
  file_size, name, ip_addr, chunks_left, chunks_total, user_agent
) VALUES (
  $1, $2, $3, $4, $5, $6
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

-- name: FileFromChunked :one
INSERT INTO files (mime_type, file_size, hash)
SELECT $2, chunked.file_size, $3 
FROM chunked
WHERE chunked.id = $1
RETURNING files.id;

-- name: UploadFromChunked :exec
INSERT INTO uploads (ip_addr, user_agent, file, name)
SELECT chunked.ip_addr, chunked.user_agent, $2, chunked.name 
FROM chunked
WHERE chunked.id = $1;

-- name: FileFromHash :one
SELECT id, mime_type FROM files
WHERE hash = $1;

-- name: FileFromId :one
SELECT mime_type, file_size FROM files
WHERE id = $1;
