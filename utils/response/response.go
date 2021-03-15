package response

import (
	"net/http"
	//"search-engine/models"

	"github.com/labstack/echo/v4"
)

func Error(c echo.Context, err error) error {
	//mega1Error, ok := err.(*models.Mega1Error)
	//if ok == true {
	//	return c.JSON(http.StatusOK, echo.Map{"code": mega1Error.Code, "message": mega1Error.Message})
	//}
	return c.JSON(http.StatusOK, echo.Map{"code": 4000, "message": err.Error()})
}

func Success(c echo.Context, msg map[string]interface{}) error {
	return c.JSON(http.StatusOK, echo.Map{"code": 2000, "data": msg})
}
