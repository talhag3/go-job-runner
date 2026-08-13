package service

import (
	"context"
	"errors"

	"github.com/talhag3/go-job-runner/internal/domain"
	"github.com/talhag3/go-job-runner/internal/repository"
	"github.com/talhag3/go-job-runner/internal/worker"
)

type JobService struct {
	repository repository.JobRepository
	workerPool *worker.WorkerPool
}

func NewJobService(repository repository.JobRepository, workerPool *worker.WorkerPool) *JobService {
	return &JobService{
		repository: repository,
		workerPool: workerPool,
	}
}

func (s *JobService) Create(ctx context.Context, job *domain.Job) error {
	if job == nil {
		return errors.New("job is required")
	}

	// Save the job.
	if err := s.repository.Create(ctx, job); err != nil {
		return err
	}

	// Send the job to workers.
	if err := s.workerPool.Submit(ctx, job); err != nil {
		return err
	}

	return nil
}

func (s *JobService) GetByID(ctx context.Context, id int64) (*domain.Job, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *JobService) GetAll(ctx context.Context) ([]*domain.Job, error) {
	return s.repository.GetAll(ctx)
}
