-- +goose Up
CREATE TYPE types AS ENUM('text', 'card', 'byte', 'auth_data');
CREATE TABLE data(
    id SERIAL primary key,
    user_id INTEGER DEFAULT NULL,
    content TEXT DEFAULT NULL,
    type types,
    is_deleted BOOLEAN DEFAULT FALSE,
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
ALTER TABLE data ADD CONSTRAINT FK_USER_TEXT_DATA FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX IDX_USER_TEXT_DATA ON data (user_Id);

-- +goose Down
ALTER TABLE data DROP CONSTRAINT FK_USER_TEXT_DATA;
DROP TABLE data;
