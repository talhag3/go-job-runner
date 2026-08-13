package executor

import (
	"context"

	"github.com/talhag3/go-job-runner/internal/domain"
)

type JobExecutor interface {
	Execute(ctx context.Context, job *domain.Job) error
}
