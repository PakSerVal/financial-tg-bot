-- +goose Up
-- +goose StatementBegin
create table selected_currency
(
    id         integer generated always as identity,
    code       text,
    user_id    bigint,
    created_at timestamp
);
create index selected_currency_user_idx on selected_currency (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table selected_currency;
-- +goose StatementEnd
