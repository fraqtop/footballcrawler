-- +goose Up
-- +goose StatementBegin
create table team (
    id serial primary key,
    title_full varchar(50) not null,
    title_short varchar(5) not null,
    unique(title_short, title_full)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table team;
-- +goose StatementEnd
