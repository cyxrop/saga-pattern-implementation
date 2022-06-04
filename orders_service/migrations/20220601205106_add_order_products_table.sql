-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_products
(
    order_id           INTEGER      NOT NULL,
    product_id         INTEGER      NOT NULL,
    number             INTEGER      NOT NULL,

    PRIMARY KEY (order_id, product_id),

    CONSTRAINT order_products_fk_order
        FOREIGN KEY (order_id)
            REFERENCES orders (id)
            ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table order_products;
-- +goose StatementEnd
