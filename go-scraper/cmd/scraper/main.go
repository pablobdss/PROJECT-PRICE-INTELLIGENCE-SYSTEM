package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/application/usecase"
	httpinfra "github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	destinationURL := os.Getenv("PRICE_DESTINATION_URL")
	if destinationURL == "" {
		log.Fatal("Erro: A variável PRICE_DESTINATION_URL é obrigatória")
	}

	sender := httpinfra.NewHTTPPriceSender(destinationURL)

	uc := usecase.NewSendPriceEventUseCase(sender)

	err = uc.Execute(
		"macbook-pro-m3",
		14999.90,
		"Amazon Brasil",
		"https://amazon.com.br/...",
		"BRL",
	)

	if err != nil {
		log.Fatalf("Erro ao enviar: %v", err)
	}

	log.Println("Evento enviado com sucesso!", destinationURL)
}
