package mongodb

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"search-engine/models"
	"search-engine/repos"
	"time"
)

type Database struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) repos.CrawlerRepository {
	return &Database{client}
}

func (db *Database) CreateCoin(ctx context.Context) error {
	var coin *models.Coin

	url := "https://coinmarketcap.com/all/views/all/"
	// Instantiate default collector
	co := colly.NewCollector()

	co.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		coin = &models.Coin{
			Name: e.ChildText(".cmc-table__column-name"),
			Symbol: e.ChildText(".cmc-table__cell--sort-by__symbol"),
			MarketCap: e.ChildText(".cmc-table__cell--sort-by__market-cap"),
			Price: e.ChildText(".cmc-table__cell--sort-by__price"),
			CirculatingSupply: e.ChildText(".cmc-table__cell--sort-by__circulating-supply"),
			Volume24h: e.ChildText(".cmc-table__cell--sort-by__volume-24-h"),
			Change1h: e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"),
			Change24h: e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h"),
			Change7d: e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		opts := options.Update().SetUpsert(true)
		filter := bson.D{{"Symbol", coin.Symbol}}
		update := bson.D{{"$set", bson.D{{"Price", ""}}}}

		_, err := db.client.Database("annashopdb").Collection("coin").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			fmt.Println("err: ", err)
		}
	})

	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	co.Visit(url)
	return nil
}
