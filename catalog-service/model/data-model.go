package model

import "time"

type ProductModel struct {
	ProductId   int       `json:"product_id" gorm:"primaryKey;autoIncrement:true"`
	Sku         string    `json:"sku"`
	Price       float64   `json:"price"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
