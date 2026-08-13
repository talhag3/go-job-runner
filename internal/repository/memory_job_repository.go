package repository

import (
	"errors"

	"github.com/talhag3/go-job-runner/internal/domain"
)

type MemoryJobRepository struct {
	jobs map[int64]*domain.Job
}

func NewMemoryJobRepository() *MemoryJobRepository {
	return &MemoryJobRepository{
		jobs: make(map[int64]*domain.Job),
	}
}

func (r *MemoryJobRepository) Create(job *domain.Job) error {

	if _, ok := r.jobs[job.ID]; ok {
		return errors.New("job already exists")
	}

	r.jobs[job.ID] = job
	return nil
}

func (r *MemoryJobRepository) GetByID(id int64) (*domain.Job, error) {

	job, ok := r.jobs[id]

	if !ok {
		return nil, errors.New("Not found")
	}

	return job, nil
}

func (r *MemoryJobRepository) GetAll() ([]*domain.Job, error) {
	jobs := []*domain.Job{}

	if len(r.jobs) == 0 {
		return nil, errors.New("No job")
	}

	for _, v := range r.jobs {
		jobs = append(jobs, v)
	}

	return jobs, nil
}
