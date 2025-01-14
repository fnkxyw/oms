-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id bigint NOT NULL,
    state text NOT NULL,
    accept_time bigint NOT NULL,
    keep_until_date timestamptz NOT NULL,
    place_date timestamptz NOT NULL,
    weight bigint NOT NULL,
    price bigint NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders;
