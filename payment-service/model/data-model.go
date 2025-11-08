package model

import "time"

type PaymentModel struct {
	PaymentId int       `json:"payment_id" gorm:"primaryKey;autoIncrement:true"`
	OrderId   int       `json:"order_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
}
