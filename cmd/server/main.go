package main

import (
	"fmt"

	"github.com/talhag3/go-job-runner/internal/domain"
	"github.com/talhag3/go-job-runner/internal/repository"
)

func main() {

	jobRepo := repository.NewMemoryJobRepository()

	job := domain.Job{ID: 1}

	err := jobRepo.Create(&job)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	jobs, err := jobRepo.GetAll()

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	for _, v := range jobs {
		fmt.Print(v.ID)
	}
}
