-- +goose Up
-- +goose StatementBegin
CREATE TABLE products
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR     NOT NULL,
    description VARCHAR     NOT NULL,
    price       INTEGER     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table products;
-- +goose StatementEnd