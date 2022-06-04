-- +goose Up
-- +goose StatementBegin
CREATE TABLE warehouses
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR     NOT NULL,
    address    VARCHAR     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table warehouses;
-- +goose StatementEnd