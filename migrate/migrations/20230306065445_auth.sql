-- +goose Up
-- +goose StatementBegin
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id            uuid PRIMARY KEY default uuid_generate_v4() NOT NULL,
    login         character varying(60)                       NOT NULL,
    password_hash character varying(80)                       NOT NULL,
    balance       decimal(10, 2)                              NOT NULL DEFAULT 0,
    created_at    timestamp        default now()              NOT NULL,
    constraint users_login_unq_ct unique (login)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
