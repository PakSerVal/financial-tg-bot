-- +goose Up
-- +goose StatementBegin
create table budget
(
    id         integer generated always as identity,
    user_id    bigint unique,
    value      bigint,
    created_at timestamp,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table budget;
-- +goose StatementEnd
