package link

import (
	"context"
	"log"
	"time"
)

// ExpiredLinksCleanupWorker periodically deletes expired links from storage.
type ExpiredLinksCleanupWorker struct {
	repo     LinkRepository
	interval time.Duration

	done    chan struct{}
	stopped chan struct{}
}

// NewExpiredLinksCleanupWorker creates a new cleanup worker.
func NewExpiredLinksCleanupWorker(repo LinkRepository, interval time.Duration) *ExpiredLinksCleanupWorker {
	return &ExpiredLinksCleanupWorker{
		repo:     repo,
		interval: interval,
		done:     make(chan struct{}),
		stopped:  make(chan struct{}),
	}
}

func (w *ExpiredLinksCleanupWorker) Start() {
	go w.run()
}

// Stop gracefully stops the worker and blocks until cleanup loop exits.
func (w *ExpiredLinksCleanupWorker) Stop() {
	close(w.done)
	<-w.stopped
}

// run is the internal background loop.
func (w *ExpiredLinksCleanupWorker) run() {
	defer close(w.stopped)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Printf("[cleanup] expired links cleanup worker started interval=%s", w.interval)

	for {
		select {
		case <-ticker.C:
			w.cleanupExpiredLinks()
		case <-w.done:
			log.Printf("[cleanup] expired links cleanup worker stopped")
			return
		}
	}
}

// cleanupExpiredLinks performs a single cleanup pass.
func (w *ExpiredLinksCleanupWorker) cleanupExpiredLinks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleted, err := w.repo.DeleteExpiredLinks(ctx)
	if err != nil {
		log.Printf("[cleanup] failed to delete expired links: %v", err)
		return
	}

	if deleted > 0 {
		log.Printf("[cleanup] deleted expired links=%d", deleted)
	}
}
