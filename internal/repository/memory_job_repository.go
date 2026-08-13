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

func (r *MemoryJobRepository) Create(ctx context.Context, job *domain.Job) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.jobs[job.ID]; ok {
		return errors.New("job already exists")
	}

	r.jobs[job.ID] = job
	return nil
}

func (r *MemoryJobRepository) GetByID(ctx context.Context, id int64) (*domain.Job, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	job, ok := r.jobs[id]

	if !ok {
		return nil, errors.New("Not found")
	}

	return job, nil
}

func (r *MemoryJobRepository) GetAll(ctx context.Context) ([]*domain.Job, error) {
	select {
	case <-ctx.Done():
		return []*domain.Job{}, ctx.Err()
	default:
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := []*domain.Job{}

	if len(r.jobs) == 0 {
		return []*domain.Job{}, nil
	}

	for _, v := range r.jobs {
		jobs = append(jobs, v)
	}

	return jobs, nil
}
