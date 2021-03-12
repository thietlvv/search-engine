package repository

import (
	"context"
	"search-engine/models"
)

type CrawlerRepository interface {
	//GetUserByMegaCode(ctx context.Context, code string) (*models.User, error)
	CreateUser(ctx context.Context, coin *models.Coin) (string, error)
}