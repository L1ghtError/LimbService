package router

import (
	"light-backend/internal/middleware"
	"light-backend/internal/ports"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(si ports.ServerInterface, r *fiber.App) {
	wrapper := ports.ServerInterfaceWrapper{
		Handler: si,
	}

	api := r.Group("/api", logger.New())

	// Auth
	auth := api.Group("/auth")

	auth.Post("/login", wrapper.Login)

	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	auth.Post("/logout", middleware.Protected([]byte(os.Getenv("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), wrapper.Logout)

	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	auth.Post("/refresh", middleware.Protected([]byte(os.Getenv("JWT_REFRESH_SECRET")), middleware.CookieTokenLookup), wrapper.Refresh)

	auth.Post("/register", wrapper.Register)

	// TODO: Implement email activation
	//auth.Get("/activate/:link", wrapper.Activate)

	// OAuth
	oauth := api.Group("/oauth")

	oauth.Get("/google", wrapper.GoogleOAuth)

	oauth.Get("/google/callback", wrapper.GoogleOAuthCb)

	// Image Related
	user := api.Group("/user")
	// TODO Refactor: instead of using Getenv we should rely on config.Get, but for now its low priority
	user.Use(middleware.Protected([]byte(os.Getenv("JWT_ACCESS_SECRET")), middleware.HeaderTokenLookup))

	user.Get("/:username", wrapper.GetUserInfo)

	user.Get("/download/image/:imageId", wrapper.DownloadImage)

	user.Post("/process/image", wrapper.ProcessImage)

	user.Post("/upload/image", wrapper.UploadImage)
}
