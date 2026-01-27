package infrastructure

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/domain/price"
)

type HTTPPriceSender struct {
	endpoint string
	client   *http.Client
}

func NewHTTPPriceSender(endpoint string) *HTTPPriceSender {
	return &HTTPPriceSender{
		endpoint: endpoint,
		client:   &http.Client{},
	}
}

func (s *HTTPPriceSender) Send(event price.PriceEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.endpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
