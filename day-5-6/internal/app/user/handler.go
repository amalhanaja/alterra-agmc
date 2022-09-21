package user

import "github.com/labstack/echo/v4"

type USerHandler interface {
	Login(c echo.Context)
	GetById(c echo.Context)
	Create(c echo.Context)
	Update(c echo.Context)
	Delete(c echo.Context)
}
