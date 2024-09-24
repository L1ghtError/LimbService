package routers

import (
	"light-backend/config"
	"light-backend/handlers"
	"light-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")
	auth.Post("/registration", handlers.Registration)
	auth.Post("/login", handlers.Login)
	auth.Post("/logout", middleware.Protected([]byte(config.Config("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), handlers.Logout)
	auth.Post("/refresh", middleware.Protected([]byte(config.Config("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), handlers.Refresh)
	auth.Get("/activate/:link", handlers.Activate)

	// OAuth
	oauth := api.Group("/oauth")
	oauth.Get("/google", handlers.Auth)
	oauth.Get("/google/callback", handlers.Callback)

	// get
	user := api.Group("/user")
	user.Use(middleware.Protected([]byte(config.Config("JWT_ACCESS_SECRET")), middleware.HeaderTokenLookup))
	user.Get("/basics", handlers.GetBasics)

	api.Get("/download/image", handlers.DownloadImage) // TODO: set under "user"
	//api.Post("/enhance/image", handlers.EnhanceImage)  // TODO: set under "user"

	user.Post("/upload/image", handlers.UploadImage)
}
