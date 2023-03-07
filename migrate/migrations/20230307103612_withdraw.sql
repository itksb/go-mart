-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdraws
(
    order_id     varchar(255)            not null,
    amount       decimal(10, 2)          not null,
    user_id      uuid                    not null,
    processed_at timestamp default now() not null,
    constraint withdraws_unique_order_number unique (order_id),
    constraint withdraws_fk_user foreign key (user_id) references users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdraws;
-- +goose StatementEnd
