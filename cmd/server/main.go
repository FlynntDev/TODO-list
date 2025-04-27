package main

import (
	"log"
	"os"

	"TODO-list/internal/handler"
	"TODO-list/internal/repository"
	"TODO-list/internal/usecase"
	"TODO-list/pkg/infrastructure"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment")
	}

	pool, err := infrastructure.NewPostgresPool()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// централизованная обработка ошибок
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	taskRepo := repository.NewTaskRepo(pool)
	taskUC := usecase.NewTaskUsecase(taskRepo)
	taskH := handler.NewTaskHandler(taskUC)
	taskH.RegisterRoutes(app)

	port := os.Getenv("APP_PORT")
	log.Printf("starting server on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
