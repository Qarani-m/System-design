CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    inventory INT NOT NULL DEFAULT 0,
    category_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


CREATE SUBSCRIPTION cqrs_sub
CONNECTION 'host=172.17.0.2 port=5432 dbname=cqrs_write user=replicator password=repl_password'
PUBLICATION cqrs_pub;


