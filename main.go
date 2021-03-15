package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"search-engine/db/main_db"

	httpHandler "search-engine/features/coinmarketcap/delivery/http"
	clientRepo "search-engine/features/coinmarketcap/repository/mongodb"
	crawlerUsecase "search-engine/features/coinmarketcap/usecase"
	//"io"
)

// Initiate web server
func main() {
	client := main_db.InitDataLayer()
	defer client.Disconnect(context.Background())
	//db := client.Database("annashopdb")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodDelete, http.MethodPut},
		MaxAge:       86400,
	}))

	httpHandler.NewHttpHandler(e, crawlerUsecase.NewUsecase(clientRepo.NewRepository(client)))

	log.Fatal(e.Start(":9090"))
}
