package price

import "time"

type PriceEvent struct {
	ProductID string
	Price     float64
	Store     string
	URL       string
	Currency  string
	Timestamp time.Time
}