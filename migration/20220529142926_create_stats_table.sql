-- +goose Up
-- +goose StatementBegin
create table stats(
    team_id int not null references team (id),
    competition_id int not null references competition (id),
    season varchar(9),
    games int not null default 0,
    points int not null default 0,
    wins int not null default 0,
    draws int not null default 0,
    losses int not null default 0,
    scored int not null default 0,
    passed int not null default 0,
    primary key (team_id, competition_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table stats
-- +goose StatementEnd
