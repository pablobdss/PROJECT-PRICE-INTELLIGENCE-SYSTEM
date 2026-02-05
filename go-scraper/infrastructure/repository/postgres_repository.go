package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pablobdss/PROJECT-PRICE-INTELLIGENCE-SYSTEM/domain/price"
)

type PostgresPriceRepository struct {
	db *sql.DB
}

func NewPostgresPriceRepository(db *sql.DB) *PostgresPriceRepository {
	return &PostgresPriceRepository{db: db}
}

func (r *PostgresPriceRepository) Save(ctx context.Context, p price.PriceEvent) error {
	query := `
		INSERT INTO prices (product_id, price, store, url, currency, scraped_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	metadata := "{}"

	_, err := r.db.ExecContext(ctx, query,
		p.ProductID,
		p.Price,
		p.Store,
		p.URL,
		p.Currency,
		p.Timestamp,
		metadata,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar no banco: %w", err)
	}

	return nil
}