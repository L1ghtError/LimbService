package router

import (
	"light-backend/internal/handlers"
	"light-backend/internal/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/registration", handlers.Registration)
	auth.Post("/login", handlers.Login)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	auth.Post("/logout", middleware.Protected([]byte(os.Getenv("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), handlers.Logout)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	auth.Post("/refresh", middleware.Protected([]byte(os.Getenv("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), handlers.Refresh)
	auth.Get("/activate/:link", handlers.Activate)

	// OAuth
	oauth := api.Group("/oauth")
	oauth.Get("/google", handlers.Auth)
	oauth.Get("/google/callback", handlers.Callback)

	// get
	user := api.Group("/user")
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	user.Use(middleware.Protected([]byte(os.Getenv("JWT_ACCESS_SECRET")), middleware.HeaderTokenLookup))
	user.Get("/basics", handlers.GetBasics)

	api.Get("/download/image/:id", handlers.DownloadImage) // TODO: remove Debug
	user.Get("/download/image/:id", handlers.DownloadImage)

	api.Post("/enhance/image", handlers.EnhanceImage) // TODO: remove Debug
	user.Post("/enhance/image", handlers.EnhanceImage)

	user.Post("/upload/image", handlers.UploadImage)
}
