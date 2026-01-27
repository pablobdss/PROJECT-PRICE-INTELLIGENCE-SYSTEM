package main

import (
	"log"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/application/usecase"
	httpinfra "github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/http"
)

func main() {
	sender := httpinfra.NewHTTPPriceSender(
		"http://localhost:8000/prices",
	)

	uc := usecase.NewSendPriceEventUseCase(sender)

	err := uc.Execute(
		"product-123",
		199.90,
		"BRL",
	)

	if err != nil {
		log.Fatal(err)
	}
}
