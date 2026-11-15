package async

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"lume_api/internal/repositories"
)

type StatsJob struct {
	UserID   string
	Username string
	Year     int
}

type StatsWorkerPool struct {
	jobs     chan StatsJob
	wg       sync.WaitGroup
	cancel   context.CancelFunc
	statsRepo *repositories.StatsRepository
	slog     *slog.Logger
}

func NewStatsWorkerPool(statsRepo *repositories.StatsRepository, logger *slog.Logger, workers, buffer int) *StatsWorkerPool {
	if workers < 1 {
		workers = 1
	}
	if buffer < 1 {
		buffer = 100
	}
	ctx, cancel := context.WithCancel(context.Background())
	p := &StatsWorkerPool{jobs: make(chan StatsJob, buffer), cancel: cancel, statsRepo: statsRepo, slog: logger}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i+1)
	}
	return p
}

func (p *StatsWorkerPool) Enqueue(job StatsJob) {
	select {
	case p.jobs <- job:
	default:
		p.slog.Warn("fila de estatísticas cheia; descartando job", "userId", job.UserID, "year", job.Year)
	}
}

func (p *StatsWorkerPool) worker(ctx context.Context, workerID int) {
	defer p.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-p.jobs:
			p.processJob(job, workerID)
		}
	}
}

func (p *StatsWorkerPool) processJob(job StatsJob, workerID int) {
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	stats, err := p.statsRepo.GetYearStats(ctx, job.Username, job.Year)
	if err != nil {
		p.slog.Error("falha ao recalcular estatísticas", "worker", workerID, "userId", job.UserID, "err", err)
		return
	}
	payload, err := json.Marshal(stats)
	if err != nil {
		p.slog.Error("falha ao serializar estatísticas", "worker", workerID, "userId", job.UserID, "err", err)
		return
	}
	if err := p.statsRepo.UpsertCachedStats(ctx, job.UserID, job.Year, payload); err != nil {
		p.slog.Error("falha ao persistir cache de estatísticas", "worker", workerID, "userId", job.UserID, "err", err)
		return
	}
	p.slog.Info("estatísticas recalculadas", "worker", workerID, "userId", job.UserID, "year", job.Year)
}

func (p *StatsWorkerPool) Shutdown() {
	p.cancel()
	p.wg.Wait()
}
