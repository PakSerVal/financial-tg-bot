-- +goose Up
-- +goose StatementBegin
create table currency_rate
(
    id         integer generated always as identity,
    code       text,
    value      bigint,
    created_at timestamp
);
create index currency_rate_code_idx on currency_rate (code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table currency_rate;
-- +goose StatementEnd
