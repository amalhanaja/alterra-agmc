package routes

import (
	"alterra-agmc-day-5-6/controllers"
	"alterra-agmc-day-5-6/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	middlewares.UseLogMiddleware(e)

	v1 := e.Group("/v1")
	v1.POST("/login", controllers.LoginUser)

	jwtMiddleware := middlewares.JWT()

	books := v1.Group("/books")
	books.POST("", controllers.CreateBook, jwtMiddleware)
	books.GET("", controllers.GetBooks)
	books.GET("/:id", controllers.GetBookById)
	books.PUT("/:id", controllers.UpdateBook, jwtMiddleware)
	books.DELETE("/:id", controllers.DeleteBook, jwtMiddleware)

	users := v1.Group("/users")
	users.POST("", controllers.CreateUser)
	users.GET("", controllers.GetUsers, jwtMiddleware)
	users.GET("/:id", controllers.GetUserById, jwtMiddleware)
	users.PUT("/:id", controllers.UpdateUser, jwtMiddleware)
	users.DELETE("/:id", controllers.DeleteUser, jwtMiddleware)

	return e
}
