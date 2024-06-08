-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS counter(
   name VARCHAR(255) PRIMARY KEY,
   value BIGINT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS counter;
-- +goose StatementEnd
