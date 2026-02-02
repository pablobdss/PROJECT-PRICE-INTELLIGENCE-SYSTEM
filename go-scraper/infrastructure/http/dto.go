package infrastructure

import (
	"time"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/domain/price"
)

type priceEventDTO struct {
	ProductID string    `json:"product_id"`
	Price     float64   `json:"price"`
	Store     string    `json:"loja"`
	URL       string    `json:"url"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

func newPriceEventDTO(e price.PriceEvent) priceEventDTO {
	return priceEventDTO{
		ProductID: e.ProductID,
		Price:     e.Price,
		Store:     e.Store,
		URL:       e.URL,
		Currency:  e.Currency,
		Timestamp: e.Timestamp,
	}
}
