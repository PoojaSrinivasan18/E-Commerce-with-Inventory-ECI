package model

import "time"

type PaymentModel struct {
	PaymentId      int       `json:"payment_id" gorm:"primaryKey;autoIncrement:true"`
	OrderId        string    `json:"order_id"`
	Amount         float64   `json:"amount"`
	Method         string    `json:"method"`
	Status         string    `json:"status"`
	Reference      string    `json:"reference"`
	IdempotencyKey string    `json:"idempotency_key" gorm:"uniqueIndex"`
	CustomerId     int       `json:"customer_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ChargeRequest represents a payment charge request
type ChargeRequest struct {
	OrderId        string  `json:"order_id" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	CustomerId     int     `json:"customer_id,omitempty"`
	Method         string  `json:"method"`
	IdempotencyKey string  `json:"idempotency_key" binding:"required"`
}

// RefundRequest represents a payment refund request
type RefundRequest struct {
	Amount float64 `json:"amount,omitempty"`
	Reason string  `json:"reason"`
}
