// main.go
package main

import (
	"log"

	handlers "github.com/aldysp34/mnc_test/rest/handler"
	"github.com/aldysp34/mnc_test/rest/middleware"
	"github.com/aldysp34/mnc_test/rest/model"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=mnc port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.TopUp{})

	e := echo.New()

	// Middleware to inject DB connection
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	// Protected routes
	protected := e.Group("/user")
	protected.Use(middleware.JWTMiddleware)
	protected.POST("/topup", handlers.TopUp)

	e.Logger.Fatal(e.Start(":8080"))
}
