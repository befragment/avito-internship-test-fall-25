-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS users_username_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS users_username_key ON users(username);
-- +goose StatementEnd

