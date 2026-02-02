package usecase

import (
	"time"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/domain/price"
)

type SendPriceEventUseCase struct {
	sender price.PriceSender
}

func NewSendPriceEventUseCase(sender price.PriceSender) *SendPriceEventUseCase {
	return &SendPriceEventUseCase{
		sender: sender,
	}
}

func (uc *SendPriceEventUseCase) Execute(
	productID string,
	value float64,
	store string, 
	url string,
	currency string,
) error {

	event := price.PriceEvent{
		ProductID: productID,
		Price:     value,
		Store:     store,
		URL:       url,
		Currency:  currency,
		Timestamp: time.Now().UTC(),
	}

	return uc.sender.Send(event)
}
