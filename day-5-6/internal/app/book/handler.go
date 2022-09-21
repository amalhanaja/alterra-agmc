package book

import "github.com/labstack/echo/v4"

type BookHandler interface {
	GetAll(c echo.Context)
	GetById(c echo.Context)
	Update(c echo.Context)
	Delete(c echo.Context)
	Create(c echo.Context)
}
