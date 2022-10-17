-- +goose Up
-- +goose StatementBegin
create table spend
(
    id         integer generated always as identity,
    price      bigint,
    user_id    bigint,
    category   text,
    created_at timestamp
);

-- Составной b-tree индекс по столбцам user_id и created_at
-- Выбрал b-tree, так как поиск происходит по операторам "=" и ">". Столбец user_id стоит первым, так как обладает большей селективностью
create index spend_user_id_created_at_idx on spend (user_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table spend;
-- +goose StatementEnd
