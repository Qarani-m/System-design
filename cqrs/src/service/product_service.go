package service

import (
	"context"
	"main/model"
	"main/repository"
)

type ProductService struct {
	writeRepo *repository.ProductWriteRepository
	readRepo  *repository.ProductReadRepository
}

 
func NewProductService(w *repository.ProductWriteRepository, r *repository.ProductReadRepository) *ProductService {
	return &ProductService{writeRepo: w, readRepo: r}
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) error {
	if p.Inventory < 0 {
		return ErrInvalidInventory
	}
	return s.writeRepo.Create(ctx, p)
}

func (s *ProductService) FindAll(ctx context.Context) ([]model.Product, error) {
	return s.readRepo.FindAll(ctx)
}

// Domain error
var ErrInvalidInventory = errorString("inventory cannot be negative")

type errorString string

func (e errorString) Error() string { return string(e) }
