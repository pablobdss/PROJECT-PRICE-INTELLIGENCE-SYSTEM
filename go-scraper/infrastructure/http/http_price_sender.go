package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/domain/price"
)

type HTTPPriceSender struct {
	endpoint string
	client   *http.Client
}

func NewHTTPPriceSender(endpoint string) *HTTPPriceSender {
	return &HTTPPriceSender{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *HTTPPriceSender) Send(event price.PriceEvent) error {
	dto := newPriceEventDTO(event)

	body, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("erro ao criar JSON: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.endpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro de conexão (timeout ou rede): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("erro na API remota: status %d", resp.StatusCode)
	}

	return nil
}
