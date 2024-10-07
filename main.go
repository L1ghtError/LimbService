package main

import (
	"fmt"
	"light-backend/amqpclient"
	"light-backend/mongoclient"
	"light-backend/routers"
	"light-backend/validation"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("example.env")
	err := mongoclient.Connect()
	if err != nil {
		fmt.Printf("MONGO %s", err.Error())
		return
	}
	err = amqpclient.Init()
	if err != nil {
		fmt.Printf("AMQP %s", err.Error())
		return
	}

	app := fiber.New(fiber.Config{
		// Global custom error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusBadRequest
			if fiberErr, ok := err.(*fiber.Error); ok {
				code = fiberErr.Code
			}
			return c.Status(code).JSON(validation.GlobalErrorHandlerResp{
				Message: err.Error(),
			})
		},
		StreamRequestBody: true, // to be able to steram media files
	})

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: "./docs/v1/swagger.yaml",
		Path:     "docs",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.1.100:5173, http://localhost:5173", // Allow all origins or specify your domains
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Content-Disposition, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Vary, Content-Length, Content-Type, Content-Disposition, ETag",
	}))

	routers.Routes(app)
	app.Listen("localhost:5266")
}
