package usecase

import (
	"context"
)

type CrawlerUsecase interface {
	CreateCoin(ctx context.Context) error
}
