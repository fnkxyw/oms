-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_state_refunded ON orders(state)
WHERE state = 'refunded';

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_state_refunded;
