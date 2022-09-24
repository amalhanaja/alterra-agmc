package app

import (
	"alterra-agmc-day-5-6/controllers"
	"alterra-agmc-day-5-6/middlewares"
	"alterra-agmc-day-5-6/pkg/app"
	"alterra-agmc-day-5-6/pkg/validator"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type restApiApp struct {
}

// OnDestroy implements app.App
func (*restApiApp) OnDestroy() {
	fmt.Println("Rest api app destroyed")
}

// OnInit implements app.App
func (*restApiApp) OnInit() error {
	panic("unimplemented")
}

// Run implements app.App
func (a *restApiApp) Run() error {
	if err := a.OnInit(); err != nil {
		return fmt.Errorf("initialization failed: %v", err)
	}
	defer a.OnDestroy()
	e := a.echo()
	e.Validator = validator.NewCustomValidator()
	return e.Start(":8080")
}

func (a *restApiApp) echo() *echo.Echo {
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

func NewRestApiApp() app.App {
	return &restApiApp{}
}
