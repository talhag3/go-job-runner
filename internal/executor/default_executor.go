package executor

import (
	"context"
	"fmt"

	"github.com/talhag3/go-job-runner/internal/domain"
)

type DefaultExecutor struct{}

func NewDefaultExecutor() *DefaultExecutor {
	return &DefaultExecutor{}
}

func (e *DefaultExecutor) Execute(ctx context.Context, job *domain.Job) error {
	fmt.Printf(
		"Executing Job %d of type %s\n",
		job.ID,
		job.Type,
	)

	return nil
}
