package repository

import (
    "context"
    "main/model"
    "encoding/json"
    "github.com/jackc/pgx/v5/pgxpool"
)

type ProductWriteRepository struct {
    db *pgxpool.Pool
}

func NewProductWriteRepository(db *pgxpool.Pool) *ProductWriteRepository {
    return &ProductWriteRepository{db: db}
}

func (r *ProductWriteRepository) Create(ctx context.Context, p *model.Product) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    err = tx.QueryRow(ctx,
        `INSERT INTO products (name, price, inventory, category_id)
         VALUES ($1,$2,$3,$4)
         RETURNING id, created_at, updated_at`,
        p.Name, p.Price, p.Inventory, p.CategoryID,
    ).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
    if err != nil {
        return err
    }

    // Insert event into outbox for syncing read DB
    payload, _ := json.Marshal(p)
    _, err = tx.Exec(ctx,
        `INSERT INTO outbox (event_type, payload) VALUES ($1,$2)`,
        "ProductCreated", payload)
    if err != nil {
        return err
    }

    return tx.Commit(ctx)
}
