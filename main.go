package main

import (
	"fmt"
	"light-backend/mongoclient"
	"light-backend/routers"
	"light-backend/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("example.env")
	err := mongoclient.Connect()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}

	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
				Success: false,
				Message: err.Error(),
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.251.102:5173, http://localhost:5173", // Allow all origins or specify your domains
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Vary, Content-Length, Content-Type, ETag",
	}))

	routers.Routes(app)
	app.Listen("localhost:5266")
}
