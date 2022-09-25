package app

import (
	"alterra-agmc-day-5-6/config"
	"alterra-agmc-day-5-6/internal/datasources"
	"alterra-agmc-day-5-6/internal/datasources/models"
	"alterra-agmc-day-5-6/internal/services"
	"alterra-agmc-day-5-6/internal/transportlayers/http/handlers"
	"alterra-agmc-day-5-6/internal/transportlayers/http/middlewares"
	"alterra-agmc-day-5-6/pkg/app"
	"alterra-agmc-day-5-6/pkg/validator"
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db, err := a.connectGormDB()
	if err != nil {
		return err
	}
	mongoDB, err := a.connectMongo(context.Background())
	if err != nil {
		return err
	}
	if mongoDB != nil {
		fmt.Println("MongoDB Connected")
	}

	// Repositories
	bookRepository := datasources.NewBookMongoDataSource(mongoDB)
	userRepository := datasources.NewUserGormDataSource(db)

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

func (a *restApiApp) connectGormDB() (*gorm.DB, error) {
	dsn := config.GetEnvOrDefault("DB_DSN", "root:password@tcp(localhost:3306)/development?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	log.Println("DB Connected")

	db.AutoMigrate(&models.UserGormModel{})
	return db, nil
}

func (a *restApiApp) connectMongo(ctx context.Context) (*mongo.Database, error) {

	clientOptions := options.Client()
	clientOptions.ApplyURI(config.GetEnvOrDefault("MONGO_URI", "mongodb://root:password@localhost:27017"))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client.Database(config.GetEnvOrDefault("MONGO_DB_NAME", "development")), nil

}

func NewRestApiApp() app.App {
	return &restApiApp{}
}
