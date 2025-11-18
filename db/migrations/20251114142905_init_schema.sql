-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS countries
(
    id SERIAL PRIMARY KEY,
    title varchar(100) NOT NULL unique,
    code varchar(2) NOT NULL,
    iso_code varchar(3)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS countries;
-- +goose StatementEnd
