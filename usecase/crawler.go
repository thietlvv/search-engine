package usecase
import (
	"context"
)

type CrawlerUsecase interface {
	processSinglePage(ctx context.Context) (string, error)
}