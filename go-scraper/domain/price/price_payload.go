package price

import "time"

type PriceEvent struct {
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}