package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/talhag3/go-job-runner/internal/executor"
	"github.com/talhag3/go-job-runner/internal/handler"
	"github.com/talhag3/go-job-runner/internal/repository"
	"github.com/talhag3/go-job-runner/internal/service"
	"github.com/talhag3/go-job-runner/internal/worker"
)

func main() {
	// Application context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --------------------------------------------------
	// Repository
	// --------------------------------------------------

	jobRepository := repository.NewMemoryJobRepository()

	// --------------------------------------------------
	// Executor
	// --------------------------------------------------

	jobExecutor := executor.NewDefaultExecutor()

	// --------------------------------------------------
	// Worker Pool
	// --------------------------------------------------

	workerPool := worker.NewWorkerPool(
		100,
		jobExecutor,
	)

	workerPool.Start(ctx, 3)

	// --------------------------------------------------
	// Service
	// --------------------------------------------------

	jobService := service.NewJobService(
		jobRepository,
		workerPool,
	)

	// --------------------------------------------------
	// HTTP Handler
	// --------------------------------------------------

	jobHandler := handler.NewJobHandler(jobService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /jobs", jobHandler.Create)
	mux.HandleFunc("GET /jobs", jobHandler.GetAll)
	mux.HandleFunc("GET /jobs/{id}", jobHandler.GetByID)

	// --------------------------------------------------
	// HTTP Server
	// --------------------------------------------------

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start HTTP server in its own goroutine.
	go func() {
		fmt.Println("HTTP server running on http://localhost:8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// --------------------------------------------------
	// OS Signals
	// --------------------------------------------------

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	fmt.Println("Job worker pool started with 3 workers.")
	fmt.Println("Press CTRL+C to shutdown.")

	// Wait for CTRL+C.
	<-signalChan

	fmt.Println("\nShutdown signal received.")

	// --------------------------------------------------
	// Cancel workers
	// --------------------------------------------------

	cancel()

	// --------------------------------------------------
	// Shutdown HTTP server
	// --------------------------------------------------

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	}

	// --------------------------------------------------
	// Wait for workers
	// --------------------------------------------------

	workerPool.Wait()

	fmt.Println("All workers stopped.")
	fmt.Println("Application stopped.")
}
