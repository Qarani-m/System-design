package repository

import (
	"main/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct{
	db *pgxpool.Pool
}


func NewProductRepository(db *pgxpool.Pool)*ProductRepository{
	return &ProductRepository{db:db}
}