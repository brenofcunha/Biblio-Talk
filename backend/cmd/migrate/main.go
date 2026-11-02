// CLI de migração independente.
//
// Uso:
//
//      go run ./cmd/migrate -direction=up    # aplica todas as pendentes
//      go run ./cmd/migrate -direction=up -steps=1
//      go run ./cmd/migrate -direction=down  # reverte todas
//
// O mesmo runner é invocado automaticamente no boot da API
// (cmd/api/main.go), então este comando é opcional.
package main

import (
        "context"
        "flag"
        "log"
        "os"
        "os/signal"
        "syscall"
        "time"

        "lume_api/internal/config"
        "lume_api/internal/database"
)

func main() {
        direction := flag.String("direction", "up", "direção da migração: up | down")
        steps := flag.Int("steps", 0, "número de migrações a aplicar (0 = todas)")
        flag.Parse()

        cfg := config.Load()
        ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
        defer stop()

        tctx, cancel := context.WithTimeout(ctx, 60*time.Second)
        defer cancel()

        pool, err := database.NewPool(tctx, cfg.DatabaseURL)
        if err != nil {
                log.Fatalf("falha ao conectar no banco: %v", err)
        }
        defer pool.Close()

        var dir database.MigrationDirection
        switch *direction {
        case "up":
                dir = database.DirectionUp
        case "down":
                dir = database.DirectionDown
        default:
                log.Fatalf("direção inválida %q (use up ou down)", *direction)
        }

        if err := database.RunMigrations(tctx, pool, os.DirFS("migrations"), dir, *steps); err != nil {
                log.Fatalf("erro ao executar migrações: %v", err)
        }
        log.Printf("migrações executadas: direction=%s steps=%d", *direction, *steps)
}
