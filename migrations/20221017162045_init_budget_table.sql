-- +goose Up
-- +goose StatementBegin
create table budget
(
    id         integer generated always as identity,
    user_id    bigint unique not null,
    value      bigint not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table budget;
-- +goose StatementEnd
