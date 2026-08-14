package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/talhag3/go-job-runner/internal/domain"
)

type MemoryJobRepository struct {
	mu   sync.RWMutex
	jobs map[int64]*domain.Job
}

func NewMemoryJobRepository() *MemoryJobRepository {
	return &MemoryJobRepository{
		jobs: make(map[int64]*domain.Job),
	}
}

func (r *MemoryJobRepository) Create(
	ctx context.Context,
	job *domain.Job,
) error {

	select {
	case <-ctx.Done():
		return ctx.Err()

	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.ID]; exists {
		return errors.New("job already exists")
	}

	r.jobs[job.ID] = job

	return nil
}

func (r *MemoryJobRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.Job, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[id]

	if !exists {
		return nil, errors.New("job not found")
	}

	return job, nil
}

func (r *MemoryJobRepository) GetAll(
	ctx context.Context,
) ([]*domain.Job, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]*domain.Job, 0, len(r.jobs))

	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}

	return jobs, nil
}
