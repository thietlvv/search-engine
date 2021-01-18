package main

import (
	"context"
	"log"
	"net/http"

	"github.com/thietlvv/search-engine/db/main_db"
	"github.com/thietlvv/search-engine/features/batdongsan.com.vn"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//"io"
)

// Initiate web server
func main() {

	batdongsan_com_vn.Start()

	main_db := main_db.InitDataLayer()
	defer main_db.Disconnect(context.Background())
	//defer main_db.Disconnect()

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
	log.Fatal(e.Start(":9090"))
}
