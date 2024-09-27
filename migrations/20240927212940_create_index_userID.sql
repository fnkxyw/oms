-- +goose Up
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- +goose Down
DROP INDEX idx_orders_user_id;
