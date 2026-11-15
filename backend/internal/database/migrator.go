package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
)

type MigrationDirection string

const (
	DirectionUp   MigrationDirection = "up"
	DirectionDown MigrationDirection = "down"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsFS fs.FS, direction MigrationDirection, steps int) error {
	if err := ensureMigrationsTable(ctx, pool); err != nil {
		return err
	}

	entries, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return err
	}

	type migrationPair struct {
		Version string
		UpPath  string
		DownPath string
	}
	pairsMap := map[string]*migrationPair{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		parts := strings.Split(name, ".")
		if len(parts) < 3 {
			continue
		}
		version := parts[0]
		dir := parts[len(parts)-2]
		if _, ok := pairsMap[version]; !ok {
			pairsMap[version] = &migrationPair{Version: version}
		}
		if dir == "up" {
			pairsMap[version].UpPath = name
		}
		if dir == "down" {
			pairsMap[version].DownPath = name
		}
	}

	pairs := make([]migrationPair, 0, len(pairsMap))
	for _, p := range pairsMap {
		pairs = append(pairs, *p)
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].Version < pairs[j].Version })

	applied, err := getAppliedMigrations(ctx, pool)
	if err != nil {
		return err
	}

	applyCount := 0
	if direction == DirectionUp {
		for _, p := range pairs {
			if applied[p.Version] {
				continue
			}
			if p.UpPath == "" {
				return fmt.Errorf("migration %s sem arquivo up", p.Version)
			}
			if err := runSingle(ctx, pool, migrationsFS, p.Version, p.UpPath, true); err != nil {
				return err
			}
			applyCount++
			if steps > 0 && applyCount >= steps {
				break
			}
		}
		return nil
	}

	for i := len(pairs) - 1; i >= 0; i-- {
		p := pairs[i]
		if !applied[p.Version] {
			continue
		}
		if p.DownPath == "" {
			return fmt.Errorf("migration %s sem arquivo down", p.Version)
		}
		if err := runSingle(ctx, pool, migrationsFS, p.Version, p.DownPath, false); err != nil {
			return err
		}
		applyCount++
		if steps > 0 && applyCount >= steps {
			break
		}
	}
	return nil
}

func ensureMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT now()
		)
	`)
	return err
}

func getAppliedMigrations(ctx context.Context, pool *pgxpool.Pool) (map[string]bool, error) {
	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

func runSingle(ctx context.Context, pool *pgxpool.Pool, migrationsFS fs.FS, version, fileName string, markApplied bool) error {
	content, err := fs.ReadFile(migrationsFS, filepath.Base(fileName))
	if err != nil {
		return err
	}
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("erro na migration %s: %w", version, err)
	}
	if markApplied {
		_, err = tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, version)
	} else {
		_, err = tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version)
	}
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

var ErrInvalidDirection = errors.New("direção de migration inválida")
