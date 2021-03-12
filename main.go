package main

import (
	"context"
	"log"
	"net/http"
	http2 "search-engine/features/coinmarketcap/delivery/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"search-engine/db/main_db"
	//"io"
)

// Initiate web server
func main() {
	client := main_db.InitDataLayer()
	defer client.Disconnect(context.Background())
	db := client.Database("annashopdb")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodDelete, http.MethodPut},
		MaxAge:       86400,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{"code": 2000, "result": "success"})
	})

	e.GET("/bds", func(c echo.Context) error {
		result := http2.Start(db)
		return c.JSON(http.StatusOK, echo.Map{"code": 2000, "result": "success", "data": result})
	})

	log.Fatal(e.Start(":9090"))
}
