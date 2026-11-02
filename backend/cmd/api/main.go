package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lume_api/internal/async"
	"lume_api/internal/config"
	"lume_api/internal/database"
	"lume_api/internal/handlers"
	"lume_api/internal/repositories"
	"lume_api/internal/router"
	"lume_api/internal/services"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("erro ao conectar no banco", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.RunMigrations(ctx, db, os.DirFS("migrations"), database.DirectionUp, 0); err != nil {
		logger.Error("erro ao executar migrations", "err", err)
		os.Exit(1)
	}

	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	logRepo := repositories.NewLogRepository(db)
	shelfRepo := repositories.NewShelfRepository(db)
	listRepo := repositories.NewListRepository(db)
	feedRepo := repositories.NewFeedRepository(db)
	statsRepo := repositories.NewStatsRepository(db)

	statsPool := async.NewStatsWorkerPool(statsRepo, logger, 4, 200)
	defer statsPool.Shutdown()

	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiration)
	usersSvc := services.NewUserService(userRepo, statsRepo)
	booksSvc := services.NewBookService(bookRepo, cfg.OpenLibraryURL, cfg.GoogleBooksURL, cfg.GoogleBooksAPIKey)
	logsSvc := services.NewLogService(db, logRepo, shelfRepo, userRepo, statsPool)
	shelfSvc := services.NewShelfService(shelfRepo)
	listSvc := services.NewListService(listRepo)
	feedSvc := services.NewFeedService(feedRepo)
	statsSvc := services.NewStatsService(statsRepo)

	h := handlers.New(cfg, db, userRepo, authSvc, usersSvc, booksSvc, logsSvc, shelfSvc, listSvc, feedSvc, statsSvc)
	app := router.Build(cfg, h, logger)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = app.ShutdownWithContext(shutdownCtx)
	}()

	logger.Info("API Lume iniciada", "port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		logger.Error("erro ao iniciar servidor", "err", err)
		os.Exit(1)
	}
}
