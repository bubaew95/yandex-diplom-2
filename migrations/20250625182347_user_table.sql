-- +goose Up
CREATE TABLE users(
    id SERIAL primary key,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email varchar(255) not null,
    password varchar(255) not null
);

-- +goose Down
DROP TABLE users;
