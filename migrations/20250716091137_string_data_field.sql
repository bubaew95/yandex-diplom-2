-- +goose Up
ALTER TABLE text_data ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;

-- +goose Down
DROP COLUMN is_deleted;
