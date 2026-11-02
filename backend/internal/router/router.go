package router

import (
        "log/slog"
        "time"

        "github.com/gofiber/fiber/v2"
        "github.com/gofiber/fiber/v2/middleware/cors"
        "github.com/gofiber/fiber/v2/middleware/limiter"
        "lume_api/internal/config"
        "lume_api/internal/handlers"
        "lume_api/internal/middleware"
)

func Build(cfg config.Config, h *handlers.Handler, logger *slog.Logger) *fiber.App {
        app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler})
        app.Use(middleware.RequestID())
        app.Use(middleware.SlogRequest(logger))
        app.Use(cors.New(cors.Config{AllowOrigins: cfg.AllowedOrigin, AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Request-Id"}))

        app.Get("/healthz", h.Healthz)
        v1 := app.Group("/api/v1")
        auth := v1.Group("/auth")
        auth.Use(limiter.New(limiter.Config{Max: 20, Expiration: time.Minute}))
        auth.Post("/register", h.Register)
        auth.Post("/login", h.Login)
        auth.Get("/me", middleware.AuthRequired(cfg.JWTSecret), h.Me)

        users := v1.Group("/users")
        users.Get("/:username", h.PublicProfile)
        users.Put("/me", middleware.AuthRequired(cfg.JWTSecret), h.UpdateMe)
        users.Post("/:id/follow", middleware.AuthRequired(cfg.JWTSecret), h.Follow)
        users.Delete("/:id/follow", middleware.AuthRequired(cfg.JWTSecret), h.Unfollow)
        users.Get("/:id/followers", h.Followers)
        users.Get("/:id/following", h.Following)
        users.Get("/:username/logs", h.UserLogs)
        users.Get("/:username/shelf", h.UserShelf)
        users.Get("/:username/stats", h.UserStats)

        books := v1.Group("/books")
        books.Get("/search", h.SearchBooks)
        books.Get("/popular", h.PopularBooks)
        books.Get("/:id", h.BookDetail)

        logs := v1.Group("/logs", middleware.AuthRequired(cfg.JWTSecret))
        logs.Post("/", h.CreateLog)
        logs.Get("/:id", h.GetLog)
        logs.Put("/:id", h.UpdateLog)
        logs.Delete("/:id", h.DeleteLog)
        logs.Post("/:id/like", h.LikeLog)
        logs.Delete("/:id/like", h.UnlikeLog)
        logs.Post("/:id/comments", h.AddComment)
        logs.Get("/:id/comments", h.ListComments)

        shelf := v1.Group("/shelf", middleware.AuthRequired(cfg.JWTSecret))
        shelf.Put("/", h.UpsertShelf)

        lists := v1.Group("/lists")
        lists.Get("/trending", h.TrendingLists)
        lists.Post("/", middleware.AuthRequired(cfg.JWTSecret), h.CreateList)
        lists.Get("/", middleware.AuthRequired(cfg.JWTSecret), h.MyLists)
        lists.Get("/:id", h.ListDetail)
        lists.Put("/:id", middleware.AuthRequired(cfg.JWTSecret), h.UpdateList)
        lists.Delete("/:id", middleware.AuthRequired(cfg.JWTSecret), h.DeleteList)
        lists.Post("/:id/items", middleware.AuthRequired(cfg.JWTSecret), h.UpsertListItem)
        lists.Delete("/:id/items", middleware.AuthRequired(cfg.JWTSecret), h.DeleteListItem)

        feed := v1.Group("/feed")
        feed.Get("/", middleware.AuthRequired(cfg.JWTSecret), h.FeedFollowing)
        feed.Get("/global", h.FeedGlobal)
        feed.Get("/reviews", h.FeedReviews)

        return app
}
