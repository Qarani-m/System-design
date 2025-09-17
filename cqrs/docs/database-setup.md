## Step 1: Clean Up 

```bash
# Stop and remove existing containers
docker stop postgres-write postgres-read
docker rm postgres-write postgres-read

```

## Step 2: Create Master (Write) Database Container

```bash
docker run --name postgres-write \
  -e POSTGRES_DB=cqrs_write \
  -e POSTGRES_USER=writeuser \
  -e POSTGRES_PASSWORD=writepass \
  -p 5432:5432 \
  -d postgres:15 \
  postgres \
  -c wal_level=logical \
  -c max_replication_slots=10 \
  -c max_wal_senders=10 \
  -c max_logical_replication_workers=4
```

## Step 3: Create Replica (Read) Database Container

```bash
docker run --name postgres-read \
  -e POSTGRES_DB=cqrs_read \
  -e POSTGRES_USER=readuser \
  -e POSTGRES_PASSWORD=readpass \
  -p 5433:5432 \
  -d postgres:15
```



## Step 4: Configure Master Database (Write DB)

**Connect to master:**
```bash
docker exec -it postgres-write psql -U writeuser -d cqrs_write
```

**Execute these commands in the master database:**

```sql
-- 1. Create replication user with proper privileges
CREATE USER replicator WITH REPLICATION LOGIN PASSWORD 'repl_password';

-- 2. Grant database connection rights
GRANT CONNECT ON DATABASE cqrs_write TO replicator;
GRANT USAGE ON SCHEMA public TO replicator;

-- 3. Create your application tables FIRST
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    inventory INT NOT NULL DEFAULT 0,
    category_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 4. Grant permissions on existing tables AND sequences
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
GRANT SELECT ON ALL SEQUENCES IN SCHEMA public TO replicator;

-- 5. Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON SEQUENCES TO replicator;

-- 6. Create publication for logical replication
CREATE PUBLICATION cqrs_pub FOR ALL TABLES;

-- 7. Verify publication
SELECT * FROM pg_publication;

-- 8. Insert some initial test data
INSERT INTO products (name, price, inventory, category_id) VALUES
('Laptop', 999.99, 10, 1),
('Mouse', 29.99, 50, 2),
('Keyboard', 79.99, 25, 2);

-- 9. Verify data exists
SELECT * FROM products;

-- Exit psql
\q
```
---
---
## Step 6: Configure Replica Database (Read DB)

**Connect to replica:**
```bash
docker exec -it postgres-read psql -U readuser -d cqrs_read
```

**Execute these commands in the replica database:**

```sql
-- 1. Create identical table structure
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    inventory INT NOT NULL DEFAULT 0,
    category_id INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 2. Get the master's IP address (we'll need this for connection)
-- Exit psql first to get IP
\q
```

**Get the master container's IP:**
```bash
*Open anoter terminal and do this*

docker inspect postgres-write | grep IPAddress
```

```sql
-- Create subscription (replace 172.17.0.2 with your actual IP)
CREATE SUBSCRIPTION cqrs_sub 
CONNECTION 'host=172.17.0.2 port=5432 dbname=cqrs_write user=replicator password=repl_password' 
PUBLICATION cqrs_pub;

-- Verify subscription was created
SELECT * FROM pg_subscription;

-- Check subscription status
SELECT subname, pid, received_lsn, latest_end_lsn FROM pg_stat_subscription;
```

## Step 7: Test Replication

**Check if initial data was replicated:**
```sql
-- Should show the 3 products from master
SELECT * FROM products;
```

**Test real-time replication - go back to master:**
```bash
docker exec -it postgres-write psql -U writeuser -d cqrs_write
```

```sql
-- Insert new data
INSERT INTO products (name, price, inventory, category_id) VALUES
('Monitor', 299.99, 8, 3);

-- Check it exists on master
SELECT * FROM products ORDER BY id;
```

**Check replica immediately:**
```bash
docker exec -it postgres-read psql -U readuser -d cqrs_read
```

```sql
-- Should show all 4 products within seconds
SELECT * FROM products ORDER BY id;
```