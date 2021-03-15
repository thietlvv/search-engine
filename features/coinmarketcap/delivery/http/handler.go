package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"search-engine/usecase"
	"search-engine/utils/response"
)

type HttpHandler struct {
	crawlerUseCase usecase.CrawlerUsecase
}

func NewHttpHandler(e *echo.Echo, uc usecase.CrawlerUsecase) {
	httpHandler := &HttpHandler{
		crawlerUseCase: uc,
	}

	e.GET("/bds", httpHandler.processSinglePage)
}

func (u *HttpHandler) processSinglePage(c echo.Context) error {

	ctx := c.Request().Context()
	err := u.crawlerUseCase.CreateCoin(ctx)
	if err != nil {
		log.Error(err)
		return response.Error(c, err)
	}
	return response.Success(c, echo.Map{"data": "ok"})
}
