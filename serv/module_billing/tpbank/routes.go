package tpbank

import (
	usecase "search-engine/usecase/billing"

	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	tpbankUsecase usecase.TpbankUsecase
}

func UseTpbank(g *echo.Group, uc usecase.TpbankUsecase) {
	httpHandler := &HttpHandler{
		tpbankUsecase: uc,
	}
	tpbank := g.Group("tpbank")

	tpbank.GET("/", httpHandler.sayHello)
}
