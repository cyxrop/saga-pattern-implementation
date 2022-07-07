-- +goose Up
-- +goose StatementBegin
CREATE TABLE reservations
(
    id           SERIAL PRIMARY KEY,
    product_id   INTEGER     NOT NULL,
    warehouse_id INTEGER     NOT NULL,
    number       INTEGER     NOT NULL,
    order_id     INTEGER     NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,

    CONSTRAINT reservations_fk_product
        FOREIGN KEY (product_id)
            REFERENCES products (id)
            ON DELETE CASCADE,

    CONSTRAINT reservations_fk_warehouse
        FOREIGN KEY (warehouse_id)
            REFERENCES warehouses (id)
            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table reservations;
-- +goose StatementEnd
