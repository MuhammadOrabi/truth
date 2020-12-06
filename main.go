package main

import (
	"log"
	"os"
	"strings"
	"truth/docs"
	"truth/handlers"
	middleware2 "truth/middleware"
	"truth/model"
	"truth/storage"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// CustomValidator ...
type CustomValidator struct {
	validator *validator.Validate
}

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// @title Truth Finder
// @version 1.0
// @description This is a Truth Finder server.

// @contact.name API Support
// @contact.url https://truth.orabi.me/support
// @contact.email muhammad@orabi.me

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	docs.SwaggerInfo.Host = os.Getenv("APP_DOMAIN")
	docs.SwaggerInfo.Schemes = []string{os.Getenv("APP_SCHEMA")}
}

func main() {
	db := storage.NewDB()
	db.AutoMigrate(&model.User{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, error=${error}, latency_human=${latency_human}, user_agent=${user_agent}\n",
	}))
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/docs/*", echoSwagger.WrapHandler)

	e.POST("/auth/register", handlers.Register)
	e.POST("/auth/login", handlers.Login)

	config := middleware.JWTConfig{
		Claims:     &model.JwtClaims{},
		SigningKey: []byte("secret"),
		BeforeFunc: func(context echo.Context) {
			token := context.Request().Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer") {
				context.Request().Header.Set("Authorization", "Bearer "+token)
			}
		},
	}

	r := e.Group("/sources", middleware.JWTWithConfig(config), middleware2.IsAdmin)
	r.GET("", handlers.ListSources)
	r.POST("", handlers.CreateSource)
	r.DELETE("/:id", handlers.DeleteSource)

	e.POST("/search", handlers.Search, middleware.JWTWithConfig(config))

	e.Logger.Fatal(e.Start(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")))
}
