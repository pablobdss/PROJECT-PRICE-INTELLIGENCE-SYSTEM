CREATE TABLE IF NOT EXISTS prices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    product_id VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    
    store VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    scraped_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    
    metadata JSONB DEFAULT '{}'::jsonb,

    CONSTRAINT check_price_positive CHECK (price >= 0),
    CONSTRAINT check_currency_code CHECK (LENGTH(currency) = 3)
);

CREATE INDEX IF NOT EXISTS idx_prices_product_id ON prices(product_id);
CREATE INDEX IF NOT EXISTS idx_prices_scraped_at ON prices(scraped_at);

CREATE INDEX IF NOT EXISTS idx_prices_metadata ON prices USING GIN (metadata);