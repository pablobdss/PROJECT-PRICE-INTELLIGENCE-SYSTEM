package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/application/usecase"
	httpinfra "github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/http"
	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/repository"
	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/infrastructure/scraper"
)

type Config struct {
	DB  *sql.DB
	URL string
}

type targetProduct struct {
	ID    string
	Store string
	URL   string
}

func main() {
	cfg := setupInfrastructure()
	defer cfg.DB.Close()

	priceRepo := repository.NewPostgresPriceRepository(cfg.DB)
	priceSender := httpinfra.NewHTTPPriceSender(cfg.URL)
	uc := usecase.NewSendPriceEventUseCase(priceRepo, priceSender)
	collyScraper := scraper.NewCollyScraper()

	runScraper(uc, collyScraper)
}

func setupInfrastructure() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: .env nao encontrado (OK em Docker)")
	}

	db := connectDB()

	webhookURL := os.Getenv("PRICE_DESTINATION_URL")
	if webhookURL == "" {
		log.Fatal("ERRO: PRICE_DESTINATION_URL obrigatoria")
	}

	return &Config{
		DB:  db,
		URL: webhookURL,
	}
}

func connectDB() *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Erro ao abrir driver: %v", err)
	}

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err == nil {
			log.Println("Banco Conectado!")
			return db
		}
		
		log.Printf("Tentando conectar ao banco (%d/5)...", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("ERRO CRITICO: Banco indisponivel: %v", err)
	return nil
}

func runScraper(uc *usecase.SendPriceEventUseCase, s *scraper.CollyScraper) {
	products := []targetProduct{
		{ID: "book-python", Store: "BooksToScrape", URL: "http://books.toscrape.com/catalogue/a-light-in-the-attic_1000/index.html"},
		{ID: "Asus-Rog", Store: "WebScraperIO", URL: "https://webscraper.io/test-sites/e-commerce/allinone/product/53"},
	}

	log.Printf("Iniciando scrap de %d produtos...", len(products))
	var wg sync.WaitGroup

	for _, p := range products {
		wg.Add(1)
		go func(product targetProduct) {
			defer wg.Done()
			processProduct(product, s, uc)
		}(p)
	}

	wg.Wait()
	log.Println("Finalizado.")
}

func processProduct(p targetProduct, s *scraper.CollyScraper, uc *usecase.SendPriceEventUseCase) {
	log.Printf("Visitando: %s", p.URL)

	data, err := s.Scrape(p.URL)
	if err != nil {
		log.Printf("ERRO SCRAPER [%s]: %v", p.ID, err)
		return
	}

	err = uc.Execute(context.Background(), p.ID, data.Price, p.Store, p.URL, data.Currency)
	if err != nil {
		log.Printf("ERRO AO SALVAR [%s]: %v", p.ID, err)
		return
	}

	log.Printf("SUCESSO [%s]: Processado e Salvo.", p.ID)
}