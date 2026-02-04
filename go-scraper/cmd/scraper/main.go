package main

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/application/usecase"
	httpinfra "github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/http"
	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/scraper"
)

type targetProduct struct {
	ID    string
	Store string
	URL   string
}

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

	collyScraper := scraper.NewCollyScraper()

	productsToScrape := []targetProduct{
		{
            ID:    "book-python",
            Store: "BooksToScrape",
            URL:   "http://books.toscrape.com/catalogue/a-light-in-the-attic_1000/index.html",
        },
		{
            ID:    "Asus-Rog",
            Store: "WebScraperIO",
            URL:   "https://webscraper.io/test-sites/e-commerce/allinone/product/53",
        },
	}

	log.Printf("Iniciando processamento de %d produtos...\n", len(productsToScrape))

	var wg sync.WaitGroup

	for _, p := range productsToScrape {
		wg.Add(1)

		go func(product targetProduct) {
			defer wg.Done()

			log.Printf(">> Visitando: %s", product.URL)

			data, err := collyScraper.Scrape(product.URL)
			if err != nil {
				log.Printf("ERRO SCRAPER [%s]: %v", product.ID, err)
				return
			}

			err = uc.Execute(
				product.ID,
				data.Price,
				product.Store,
				product.URL,
				data.Currency,
			)

			if err != nil {
				log.Printf("ERRO [%s]: %v", product.ID, err)
			} 

			log.Printf("SUCESSO [%s]: Enviado.", product.ID)
			
		}(p)
	}

	log.Println("Aguardando todas as goroutines...")
	wg.Wait()

	log.Println("Processamento finalizado com sucesso.")
}
