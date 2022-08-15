-- +goose Up
-- +goose StatementBegin
create table competition (
    id serial primary key,
    title varchar(100) not null unique,
    uri_table varchar(255) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table competition
-- +goose StatementEnd
