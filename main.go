package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
	"log"
	"main.go/config/migration"
	routes "main.go/route"
	"os"
	"time"
)

func main() {
	env := godotenv.Load()
	migrate := os.Getenv("MIGRATE")
	if migrate == "TRUE" {
		migration.Migrate()
	}

	log.Println("Run System")

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})

	app.Use(cors.New())

	app.Use(
		logger.New(logger.Config{
			Format: "${time} | [${ip}]:${port} | ${host}${path} | ${status} - ${method}\n",
		}),
		limiter.New(limiter.Config{
			Max:        1000,
			Expiration: 30 * time.Second,
		}),
	)
	app.Get("/Log", monitor.New(monitor.Config{Title: "Monitoring Webservice CMS", Refresh: 3 * time.Second}))
	routes.Init(app)

	if env != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	listen := app.Listen(":" + port)
	if listen != nil {
		log.Println("Fail to listen go fiber server")
		os.Exit(1000)
	}
}
