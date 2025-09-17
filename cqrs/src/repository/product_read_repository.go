package repository

import (
    "context"
    "main/model"
    "github.com/jackc/pgx/v5/pgxpool"
)

type ProductReadRepository struct {
    db *pgxpool.Pool
}

func NewProductReadRepository(db *pgxpool.Pool) *ProductReadRepository {
    return &ProductReadRepository{db: db}
}

func (r *ProductReadRepository) FindAll(ctx context.Context) ([]model.Product, error) {
    rows, err := r.db.Query(ctx,
        `SELECT id, name, price, inventory, category_id, created_at, updated_at
         FROM products ORDER BY id DESC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []model.Product
    for rows.Next() {
        var p model.Product
        err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Inventory,
            &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
        if err != nil {
            return nil, err
        }
        products = append(products, p)
    }
    return products, nil
}
