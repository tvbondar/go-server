\c orders_db
CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(50),
    entry VARCHAR(10),
    locale VARCHAR(10),
    internal_signature VARCHAR(50),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10),
    delivery JSONB,
    payment JSONB
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid),
    chrt_id BIGINT,
    track_number VARCHAR(50),
    price INT,
    rid VARCHAR(50),
    name VARCHAR(100),
    sale INT,
    size VARCHAR(10),
    total_price INT,
    nm_id BIGINT,
    brand VARCHAR(100),
    status INT
);

CREATE INDEX idx_orders_order_uid ON orders(order_uid);