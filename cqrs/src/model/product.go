package model

import "time"

type Product struct {
	ID int `json:id`
	Name string `json:name`  
	Price float64 `json:price`
	Inventory int `json:inventory`
	CategoryID int `json:category_id`
	CreatedAt time.Time `json:create_at`
	UpdatedAt time.Time `json:update_at`
}