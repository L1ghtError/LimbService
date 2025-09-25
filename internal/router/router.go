package router

import (
	"light-backend/internal/handlers"
	"light-backend/internal/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(handler handlers.HttpServer, r *fiber.App) {
	api := r.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/registration", handler.Registration)
	auth.Post("/login", handler.Login)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	auth.Post("/logout", middleware.Protected([]byte(os.Getenv("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), handler.Logout)
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	auth.Post("/refresh", middleware.Protected([]byte(os.Getenv("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), handler.Refresh)
	auth.Get("/activate/:link", handler.Activate)

	// OAuth
	oauth := api.Group("/oauth")
	oauth.Get("/google", handler.Auth)
	oauth.Get("/google/callback", handler.Callback)

	// get
	user := api.Group("/user")
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	user.Use(middleware.Protected([]byte(os.Getenv("JWT_ACCESS_SECRET")), middleware.HeaderTokenLookup))
	user.Get("/basics", handler.GetBasics)

	api.Get("/download/image/:id", handler.DownloadImage) // TODO: remove Debug
	user.Get("/download/image/:id", handler.DownloadImage)

	api.Post("/enhance/image", handler.EnhanceImage) // TODO: remove Debug
	user.Post("/enhance/image", handler.EnhanceImage)

	user.Post("/upload/image", handler.UploadImage)
}
