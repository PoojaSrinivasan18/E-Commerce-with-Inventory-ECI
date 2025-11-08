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
