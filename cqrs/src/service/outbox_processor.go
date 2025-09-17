package service

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "main/model"
    "github.com/jackc/pgx/v5/pgxpool"
)

type OutboxProcessor struct {
    writeDB *pgxpool.Pool
    readDB  *pgxpool.Pool
}

func NewOutboxProcessor(w, r *pgxpool.Pool) *OutboxProcessor {
    return &OutboxProcessor{writeDB: w, readDB: r}
}

func (p *OutboxProcessor) Start() {
    ticker := time.NewTicker(1 * time.Second)
    ctx := context.Background()

    go func() {
        for range ticker.C {
            rows, err := p.writeDB.Query(ctx,
                `SELECT id, event_type, payload 
                 FROM outbox 
                 WHERE processed_at IS NULL 
                 ORDER BY id 
                 LIMIT 10`)
            if err != nil {
                log.Println("outbox query error:", err)
                continue
            }

            var ids []int64
            for rows.Next() {
                var id int64
                var eventType string
                var payload []byte
                if err := rows.Scan(&id, &eventType, &payload); err != nil {
                    log.Println("scan error:", err)
                    continue
                }

                switch eventType {
                case "ProductCreated":
                    var prod model.Product
                    if err := json.Unmarshal(payload, &prod); err != nil {
                        log.Println("unmarshal error:", err)
                        continue
                    }
                    _, err := p.readDB.Exec(ctx,
                        `INSERT INTO products (id, name, price, inventory, category_id, created_at, updated_at)
                         VALUES ($1,$2,$3,$4,$5,$6,$7)
                         ON CONFLICT (id) DO NOTHING`,
                        prod.ID, prod.Name, prod.Price, prod.Inventory,
                        prod.CategoryID, prod.CreatedAt, prod.UpdatedAt)
                    if err != nil {
                        log.Println("read db insert error:", err)
                        continue
                    }
                }
                ids = append(ids, id)
            }
            rows.Close()

            // ✅ Update processed_at instead of processed
            for _, id := range ids {
                _, err := p.writeDB.Exec(ctx, 
                    `UPDATE outbox SET processed_at = NOW() WHERE id = $1`, id)
                if err != nil {
                    log.Println("mark processed error:", err)
                }
            }
        }
    }()
}
