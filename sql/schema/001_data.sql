-- +goose Up
CREATE TABLE data (
	id INTEGER PRIMARY KEY,
	url TEXT UNIQUE NOT NULL,
	content TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

-- +goose Down
DROP TABLE data;
