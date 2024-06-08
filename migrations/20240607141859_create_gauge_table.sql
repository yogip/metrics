-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gauge(
   name VARCHAR(255) PRIMARY KEY,
   value DOUBLE PRECISION NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gauge;
-- +goose StatementEnd
