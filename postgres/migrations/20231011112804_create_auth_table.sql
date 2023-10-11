-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth (
    id serial primary key,
    name text not null,
    email text,
    role integer,
    created_at timestamp not null default now(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth;
-- +goose StatementEnd
