-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_id ON orders(user_id);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_orders_user_id;
