-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
                        order_uid VARCHAR(255) PRIMARY KEY,
                        track_number VARCHAR(255) NOT NULL,
                        entry VARCHAR(50) NOT NULL,
                        locale VARCHAR(10) NOT NULL,
                        internal_signature VARCHAR(255),
                        customer_id VARCHAR(255) NOT NULL,
                        delivery_service VARCHAR(100) NOT NULL,
                        shardkey VARCHAR(10) NOT NULL,
                        sm_id INTEGER NOT NULL,
                        date_created TIMESTAMP WITH TIME ZONE NOT NULL,
                        oof_shard VARCHAR(10) NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE deliveries (
                            id SERIAL PRIMARY KEY,
                            order_uid VARCHAR(255) NOT NULL UNIQUE REFERENCES orders(order_uid) ON DELETE CASCADE,
                            name VARCHAR(255) NOT NULL,
                            phone VARCHAR(50) NOT NULL,
                            zip VARCHAR(20) NOT NULL,
                            city VARCHAR(100) NOT NULL,
                            address TEXT NOT NULL,
                            region VARCHAR(100) NOT NULL,
                            email VARCHAR(255) NOT NULL,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
                          id SERIAL PRIMARY KEY,
                          order_uid VARCHAR(255) NOT NULL UNIQUE REFERENCES orders(order_uid) ON DELETE CASCADE,
                          transaction VARCHAR(255) NOT NULL,
                          request_id VARCHAR(255),
                          currency VARCHAR(10) NOT NULL,
                          provider VARCHAR(100) NOT NULL,
                          amount NUMERIC(12, 2) NOT NULL,
                          payment_dt BIGINT NOT NULL,
                          bank VARCHAR(100) NOT NULL,
                          delivery_cost NUMERIC(12, 2) NOT NULL,
                          goods_total INTEGER NOT NULL,
                          custom_fee INTEGER NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                       chrt_id BIGINT NOT NULL,
                       track_number VARCHAR(255) NOT NULL,
                       price NUMERIC(12, 2) NOT NULL,
                       rid VARCHAR(255) NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       sale INTEGER NOT NULL,
                       size VARCHAR(50) NOT NULL,
                       total_price NUMERIC(12, 2) NOT NULL,
                       nm_id BIGINT NOT NULL,
                       brand VARCHAR(255) NOT NULL,
                       status INTEGER NOT NULL,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS deliveries;
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
