-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id                 SERIAL PRIMARY KEY,
    warehouse_id       INTEGER      NOT NULL,
    status             SMALLINT     NOT NULL,
    status_description VARCHAR(256) NOT NULL,
    created_at         TIMESTAMPTZ  NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table orders;
-- +goose StatementEnd
