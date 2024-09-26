-- +goose Up
CREATE TABLE orders (
                        id serial primary key ,
                        user_id integer not null ,
                        state varchar(50) not null,
                        accept_time bigint,
                        keep_until_date timestamp,
                        place_date timestamp,
                        weight integer,
                        price integer
);
-- +goose Down
drop table if exists orders;