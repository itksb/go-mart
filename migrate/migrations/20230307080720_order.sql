-- +goose Up
-- +goose StatementBegin
create table orders
(
    id          varchar(255) primary key     not null,
    status      varchar(10)                  not null,
    accrual     decimal(10, 2) default 0,
    user_id     uuid                         not null,
    uploaded_at timestamp      default now() not null,
    constraint orders_user_fk foreign key (user_id) references users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
