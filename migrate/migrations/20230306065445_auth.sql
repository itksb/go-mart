-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE  users (
    id uuid primary key default uuid_generate_v4() not null,
    login character varying(100) not null,
    password_hash character varying(100) not null,
    balance decimal(10, 2) default 0,
    created_at timestamp default now() not null,
    constraint users_unique_login unique (login)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
