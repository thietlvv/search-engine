package tpbank

import (
	"search-engine/utils/response"

	"github.com/labstack/echo/v4"
)

// var HttpHandler urls.HttpHandler

func (u *HttpHandler) sayHello(c echo.Context) error {

	// ctx := c.Request().Context()
	// err := u.crawlerUseCase.CreateCoin(ctx)
	// if err != nil {
	// 	log.Error(err)
	// 	return response.Error(c, err)
	// }
	return response.Success(c, echo.Map{"data": "ok"})
}
