-- +goose Up
-- +goose StatementBegin
CREATE TABLE invoices
(
    id         SERIAL PRIMARY KEY,
    order_id   INTEGER     NOT NULL,
    status     SMALLINT    NOT NULL,
    amount     INTEGER     NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE invoices;
-- +goose StatementEnd
