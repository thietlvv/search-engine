package http

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
	//"go.mongodb.org/mongo-driver/mongo"
	"search-engine/models"
	"search-engine/usecase"

)

type HttpHandler struct {
	crawlerUseCase usecase.CrawlerUsecase
}

func NewHttpHandler(e *echo.Echo, uc usecase.CrawlerUsecase) {
	httpHandler := &HttpHandler{
		crawler_usecase: uc,
	}

	e.GET("/bds", httpHandler.processSinglePage)
}

func (u *HttpHandler) processSinglePage(c echo.Context) error {
	var coin *models.Coin

	url := "https://coinmarketcap.com/all/views/all/"
	// Instantiate default collector
	co := colly.NewCollector()

	co.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		coin := Coin{
			Name: e.ChildText(".cmc-table__column-name"),
			Symbol: e.ChildText(".cmc-table__cell--sort-by__symbol"),
			MarketCap: e.ChildText(".cmc-table__cell--sort-by__market-cap"),
			Price: e.ChildText(".cmc-table__cell--sort-by__price"),
			CirculatingSupply: e.ChildText(".cmc-table__cell--sort-by__circulating-supply"),
			Volume24h: e.ChildText(".cmc-table__cell--sort-by__volume-24-h"),
			Change1h: e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h"),
			Change24h: e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h"),
			Change7d: e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d"),
		}
	})

	co.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	co.Visit(url)
	return coins
}