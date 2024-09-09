package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/kritpi/go-task-management/configs"
	"github.com/kritpi/go-task-management/domain/usecases"
	"github.com/kritpi/go-task-management/internal/adapters/mysql"
	"github.com/kritpi/go-task-management/internal/adapters/rest"
	"github.com/kritpi/go-task-management/middlewares"
)

func main() {
	app := fiber.New()

	ctx := context.Background()

	cfg := configs.NewConfig()

	db, err := sqlx.ConnectContext(ctx, "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userRepo := mysql.NewUserMySQLRepository(db)
	userService := usecases.NewUserService(userRepo, cfg)
	userHandler := rest.NewUserHandler(userService)

	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	app.Use(middlewares.JwtMiddleware(cfg.JWTSecret))

	if err := app.Listen(":9000"); err != nil {
		log.Fatal(err)
	}
}