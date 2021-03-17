package useragent

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	ua "github.com/mileusna/useragent"
)

func ContextWithValue(c echo.Context) context.Context {
	userAgent := c.Request().Header.Get("User-Agent")
	ua := ua.Parse(userAgent)
	fmt.Printf("User-Agent %+v", ua)
	ctx := context.WithValue(c.Request().Context(), "User-Agent", ua)
	return ctx
}
