package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/talhag3/go-job-runner/internal/domain"
)

type DefaultExecutor struct{}

func NewDefaultExecutor() *DefaultExecutor {
	return &DefaultExecutor{}
}

func (e *DefaultExecutor) Execute(
	ctx context.Context,
	job *domain.Job,
) error {

	fmt.Printf(
		"Executing Job [%d] type [%s]\n",
		job.ID,
		job.Type,
	)

	// Simulate some work.
	select {
	case <-time.After(2 * time.Second):
		fmt.Printf("Job [%d] execution finished\n", job.ID)
		return nil

	case <-ctx.Done():
		fmt.Printf("Job [%d] execution cancelled\n", job.ID)
		return ctx.Err()
	}
}
