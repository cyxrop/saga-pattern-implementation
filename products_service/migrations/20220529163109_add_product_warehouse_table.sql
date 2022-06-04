-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_warehouse
(
    product_id   INTEGER NOT NULL,
    warehouse_id INTEGER NOT NULL,
    number       INTEGER NOT NULL,

    PRIMARY KEY (product_id, warehouse_id),

    CONSTRAINT product_warehouse_fk_product
        FOREIGN KEY (product_id)
            REFERENCES products (id)
            ON DELETE CASCADE,

    CONSTRAINT product_warehouse_fk_warehouse
        FOREIGN KEY (warehouse_id)
            REFERENCES warehouses (id)
            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table product_warehouse;
-- +goose StatementEnd
