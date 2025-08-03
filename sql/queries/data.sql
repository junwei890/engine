-- name: InsertData :one
INSERT OR REPLACE INTO data (url, content, created_at, updated_at) VALUES (
	?,
	?,
	?,
	?
) RETURNING url;
