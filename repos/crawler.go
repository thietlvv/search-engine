package repos

import (
	"context"
)

type CrawlerRepository interface {
	CreateCoin(ctx context.Context) error
}
