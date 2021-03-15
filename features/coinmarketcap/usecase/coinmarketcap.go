package usecase

import (
	"context"
	"search-engine/repos"
	"search-engine/usecase"
)

type UsecaseRepo struct {
	crawler_repo repos.CrawlerRepository
}

func NewUsecase(crawlerRepo repos.CrawlerRepository) usecase.CrawlerUsecase {
	return &UsecaseRepo{
		crawler_repo: crawlerRepo,
	}
}

func (u *UsecaseRepo) CreateCoin(ctx context.Context) error {

	return u.crawler_repo.CreateCoin(ctx)
}
