-- +goose Up
-- +goose StatementBegin
ALTER TABLE chat ADD COLUMN public boolean;
ALTER TABLE chat ADD COLUMN name varchar;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chat DROP COLUMN public;
ALTER TABLE chat DROP COLUMN name;
-- +goose StatementEnd
