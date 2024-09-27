-- +goose Up
CREATE INDEX idx_orders_state_refunded ON orders(state)
    WHERE state = 'refunded';


-- +goose Down
DROP INDEX idx_orders_state_refunded;