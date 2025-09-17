CREATE USER replicator WITH REPLICATION LOGIN PASSWORD 'repl_password';

GRANT CONNECT ON DATABASE cqrs_write TO replicator;
GRANT USAGE ON SCHEMA public TO replicator;

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    inventory INT NOT NULL DEFAULT 0,
    category_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS outbox (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP NULL
);

GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO replicator;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON SEQUENCES TO replicator;

CREATE PUBLICATION cqrs_pub FOR ALL TABLES;

INSERT INTO products (name, price, inventory, category_id) VALUES
('Laptop', 999.99, 10, 1),
('Mouse', 29.99, 50, 2),
('Keyboard', 79.99, 25, 2);
