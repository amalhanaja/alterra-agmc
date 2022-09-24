package app

import (
	"alterra-agmc-day-5-6/config"
	"alterra-agmc-day-5-6/internal/datasources"
	"alterra-agmc-day-5-6/internal/services"
	"alterra-agmc-day-5-6/internal/transportlayers/http/handlers"
	"alterra-agmc-day-5-6/internal/transportlayers/http/middlewares"
	"alterra-agmc-day-5-6/pkg/app"
	"alterra-agmc-day-5-6/pkg/validator"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type restApiApp struct {
	userHandler handlers.UserHandler
	bookHandler handlers.BookHandler
}

// OnDestroy implements app.App
func (*restApiApp) OnDestroy() {
	fmt.Println("Rest api app destroyed")
}

// OnInit implements app.App
func (a *restApiApp) OnInit() error {
	config.InitDB()

	// Repositories
	bookRepository := datasources.NewBookInMemoryDataSource()
	userRepository := datasources.NewUserGormDataSource(config.DB)

	// Services
	bookService := services.NewBookService(bookRepository)
	userService := services.NewUserService(userRepository)

	// Handlers
	a.bookHandler = handlers.NewBookHandler(bookService)
	a.userHandler = handlers.NewUserHandler(userService)

	return nil
}

// Run implements app.App
func (a *restApiApp) Run() error {
	if err := a.OnInit(); err != nil {
		return fmt.Errorf("initialization failed: %v", err)
	}
	defer a.OnDestroy()
	return a.echo().Start(":8080")
}

func (a *restApiApp) echo() *echo.Echo {
	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	e.Pre(middleware.RemoveTrailingSlash())
	middlewares.UseLogMiddleware(e)

	v1 := e.Group("/v1")
	v1.POST("/login", a.userHandler.Login)

	jwtMiddleware := middlewares.JWT()

	books := v1.Group("/books")
	books.POST("", a.bookHandler.Create, jwtMiddleware)
	books.GET("", a.bookHandler.GetAll)
	books.GET("/:id", a.bookHandler.GetByID)
	books.PUT("/:id", a.bookHandler.Update, jwtMiddleware)
	books.DELETE("/:id", a.bookHandler.Delete, jwtMiddleware)

	users := v1.Group("/users")
	users.POST("", a.userHandler.Create)
	users.GET("", a.userHandler.GetAll, jwtMiddleware)
	users.GET("/:id", a.userHandler.GetByID, jwtMiddleware)
	users.PUT("/:id", a.userHandler.Update, jwtMiddleware)
	users.DELETE("/:id", a.userHandler.Delete, jwtMiddleware)

	return e
}

func NewRestApiApp() app.App {
	return &restApiApp{}
}
