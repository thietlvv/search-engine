package urls

import (
	tpbank "search-engine/serv/module_billing"

	"github.com/labstack/echo/v4"
)

func InitUrlsBilling(e *echo.Echo) {

	billing := e.Group("/billing")

	tpbank.UseTpbank(billing, tpbank.sayHello)

}
