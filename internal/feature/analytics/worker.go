package analytics

import "log"

type ClickEventWorker struct {
	service *AnalyticsService
	done    chan struct{}
	stopped chan struct{}
}

func NewClickEventWorker(service *AnalyticsService) *ClickEventWorker {
	return &ClickEventWorker{
		service: service,
		done:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

// Start runs background analytics worker.
// Consumes click events asynchronously: writes them to storage.
// Performs graceful drain of the remaining queue on shutdown.
func (worker *ClickEventWorker) Start() {
	log.Println("[ClickEventWorker] worker started")
	go worker.run()
}

func (worker *ClickEventWorker) run() {
	defer close(worker.stopped)

	for {
		select {
		case ev := <-worker.service.clickEventChan:
			worker.service.handleClickEvent(ev)
		case <-worker.done:
			worker.drain()
			return
		}
	}
}

// Stop gracefully shuts down.
// Drains all remaining queued events, and blocks until the worker has fully completed.
func (worker *ClickEventWorker) Stop() {
	log.Println("[ClickEventWorker] worker shutdown initiated")
	close(worker.done)
	<-worker.stopped
	log.Println("[ClickEventWorker] worker stopped")
}

// Drain remaining events on shutdown
func (worker *ClickEventWorker) drain() {
	for {
		select {
		case ev := <-worker.service.clickEventChan:
			worker.service.handleClickEvent(ev)
		default:
			log.Println("[ClickEventWorker] Completed all clickEvents")
			return
		}
	}
}
