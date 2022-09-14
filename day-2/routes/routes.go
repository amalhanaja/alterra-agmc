package routes

import (
	"alterra-agmc-day-2/controllers"

	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	v1 := e.Group("/v1")

	books := v1.Group("/books")
	books.POST("", controllers.CreateBook)
	books.GET("", controllers.GetBooks)
	books.GET("/:id", controllers.GetBookById)
	books.PUT("/:id", controllers.UpdateBook)
	books.DELETE("/:id", controllers.DeleteBook)

	return e
}
