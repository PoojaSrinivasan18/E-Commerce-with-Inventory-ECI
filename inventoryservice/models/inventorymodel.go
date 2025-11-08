package models

import "time"

type InventoryModel struct {
	InventoryId int       `json:"inventory_id" gorm:"primaryKey;autoIncrement:true"`
	ProductId   int       `json:"product_id"`
	WareHouse   string    `json:"warehouse"`
	OnHand      int       `json:"onhand"`
	Reserved    int       `json:"reserved"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReservationRequest represents a request to reserve inventory
type ReservationRequest struct {
	ProductId      int    `json:"product_id" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	Warehouse      string `json:"warehouse,omitempty"`
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	OrderId        string `json:"order_id" binding:"required"`
}

// ReservationRecord tracks individual reservations with TTL
type ReservationRecord struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement:true"`
	ProductId      int       `json:"product_id"`
	Warehouse      string    `json:"warehouse"`
	Quantity       int       `json:"quantity"`
	OrderId        string    `json:"order_id"`
	IdempotencyKey string    `json:"idempotency_key" gorm:"uniqueIndex"`
	Status         string    `json:"status"` // RESERVED, SHIPPED, RELEASED, EXPIRED
	ReservedAt     time.Time `json:"reserved_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ReleaseRequest represents a request to release reserved inventory
type ReleaseRequest struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	OrderId        string `json:"order_id" binding:"required"`
}

// ShipRequest represents a request to ship reserved inventory
type ShipRequest struct {
	IdempotencyKey string `json:"idempotency_key" binding:"required"`
	OrderId        string `json:"order_id" binding:"required"`
}
