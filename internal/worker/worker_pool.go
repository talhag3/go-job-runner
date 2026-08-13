package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/talhag3/go-job-runner/internal/domain"
	"github.com/talhag3/go-job-runner/internal/executor"
)

type WorkerPool struct {
	jobs     chan *domain.Job
	executor executor.JobExecutor
	wg       sync.WaitGroup
}

func NewWorkerPool(queueSize int, executor executor.JobExecutor) *WorkerPool {
	return &WorkerPool{
		jobs:     make(chan *domain.Job, queueSize),
		executor: executor,
	}
}

func (p *WorkerPool) Start(ctx context.Context, workerCount int) {
	p.wg.Add(workerCount)

	for i := 1; i <= workerCount; i++ {
		go p.worker(ctx, i)
	}
}

func (p *WorkerPool) Submit(ctx context.Context, job *domain.Job) error {
	select {
	case p.jobs <- job:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *WorkerPool) worker(ctx context.Context, id int) {
	defer p.wg.Done()

	for {
		select {
		case job := <-p.jobs:
			fmt.Printf(
				"Worker %d received Job %d\n",
				id,
				job.ID,
			)

			if err := p.executor.Execute(ctx, job); err != nil {
				fmt.Printf(
					"Worker %d failed Job %d: %v\n",
					id,
					job.ID,
					err,
				)
			}

		case <-ctx.Done():
			fmt.Printf("Worker %d stopping\n", id)
			return
		}
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}
