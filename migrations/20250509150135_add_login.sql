-- +goose Up
-- +goose StatementBegin
ALTER TABLE message ADD COLUMN login varchar(36);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE message DROP COLUMN login;
-- +goose StatementEnd
