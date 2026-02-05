package price

import "context"

type Repository interface {
	Save(ctx context.Context, event PriceEvent) error
}