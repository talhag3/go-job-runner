package repository

import "github.com/talhag3/go-job-runner/internal/domain"

type JobRepository interface {
	Create(job *domain.Job) error
	GetByID(id int64) (*domain.Job, error)
	GetAll() ([]*domain.Job, error)
}
