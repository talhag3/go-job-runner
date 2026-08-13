package repository

import (
	"context"

	"github.com/talhag3/go-job-runner/internal/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	GetByID(ctx context.Context, id int64) (*domain.Job, error)
	GetAll(ctx context.Context) ([]*domain.Job, error)
}
