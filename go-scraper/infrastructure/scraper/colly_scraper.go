package scraper

import (
	"fmt"
	"log" 
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type ScrapedData struct {
	Name     string
	Price    float64
	Currency string
}

type CollyScraper struct {
	collector *colly.Collector
}

func NewCollyScraper() *CollyScraper {
	c := colly.NewCollector(

		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(10 * time.Second)

	return &CollyScraper{collector: c}
}

func (s *CollyScraper) Scrape(url string) (*ScrapedData, error) {
	var data ScrapedData
	var err error

	c := s.collector.Clone()
	parserFound := false

	
	c.OnResponse(func(r *colly.Response) {
		log.Printf("[DEBUG] Status: %d | URL: %s", r.StatusCode, r.Request.URL.String())
	})

	
	if strings.Contains(url, "toscrape") {
		parserFound = true
		
		c.OnHTML("div.product_main h1", func(e *colly.HTMLElement) {
			data.Name = strings.TrimSpace(e.Text)
		})
		
		c.OnHTML("p.price_color", func(e *colly.HTMLElement) {
			
			priceStr := strings.ReplaceAll(e.Text, "£", "")
			priceStr = strings.TrimSpace(priceStr)
			data.Price, data.Currency = parsePrice(e.Text)
		})
	}

	if strings.Contains(url, "webscraper.io") {
		parserFound = true
		
		
		c.OnHTML(".caption h4:not(.price)", func(e *colly.HTMLElement) {
			
			data.Name = strings.TrimSpace(e.Text)
		})

		
		c.OnHTML(".caption h4.price", func(e *colly.HTMLElement) {
			
			data.Price, data.Currency = parsePrice(e.Text)
		})
	}

	
	if !parserFound {
		c.OnHTML("title", func(e *colly.HTMLElement) {
			data.Name = strings.TrimSpace(e.Text)
		})
	}

	
	c.OnError(func(r *colly.Response, e error) {
		err = fmt.Errorf("falha: %v (Status: %d)", e, r.StatusCode)
	})

	visitErr := c.Visit(url)
	if visitErr != nil {return nil, visitErr}

	if err != nil {return nil, err}

	if data.Name == "" {
		return nil, fmt.Errorf("nome não encontrado (Provável Bloqueio/Redirect)")
	}

	return &data, nil
}
