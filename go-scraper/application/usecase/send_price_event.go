package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/domain/price"
)

type SendPriceEventUseCase struct {
	repo   price.Repository
	sender price.PriceSender 
}

func NewSendPriceEventUseCase(repo price.Repository, sender price.PriceSender) *SendPriceEventUseCase {
	return &SendPriceEventUseCase{
		repo:   repo,
		sender: sender,
	}
}

func (uc *SendPriceEventUseCase) Execute(
	ctx context.Context,
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

	if err := uc.repo.Save(ctx, event); err != nil {
		return fmt.Errorf("falha ao persistir preço: %w", err)
	}
	log.Printf("Preço salvo no Banco: %s - %.2f", productID, value)

	if err := uc.sender.Send(event); err != nil {
		log.Printf("Falha ao enviar webhook: %v", err)
	}

	return nil
}