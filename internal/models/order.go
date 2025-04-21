package models

import (
	"time"
)

type Cart struct {
	UserID    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ProductID   string  `json:"product_id" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}

type Order struct {
	OrderNumber string     `json:"order_numder" validate:"required"`
	UserID      string     `json:"user_id" validate:"required"`
	Items       []CartItem `json:"items" validate:"required"`
	CreatedAt   time.Time  `json:"created_at" validate:"required"`
	Status      string     `json:"status" validate:"required"`
	LastStatus  string     `json:"last_status"`
	TotalPrice  float64    `json:"total_price" validate:"required"`
}
