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
	// ------------------------------------------------
	// Application context
	// ------------------------------------------------

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ------------------------------------------------
	// Dependencies
	// ------------------------------------------------

	jobRepository := repository.NewMemoryJobRepository()

	jobExecutor := executor.NewDefaultExecutor()

	workerPool := worker.NewWorkerPool(
		100,
		jobExecutor,
	)

	// Start 3 workers.
	workerPool.Start(ctx, 3)

	// ------------------------------------------------
	// Service
	// ------------------------------------------------

	jobService := service.NewJobService(
		jobRepository,
		workerPool,
	)

	// ------------------------------------------------
	// HTTP Handler
	// ------------------------------------------------

	jobHandler := handler.NewJobHandler(jobService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /jobs", jobHandler.Create)
	mux.HandleFunc("GET /jobs", jobHandler.GetAll)
	mux.HandleFunc("GET /jobs/{id}", jobHandler.GetByID)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// ------------------------------------------------
	// Start HTTP server
	// ------------------------------------------------

	go func() {
		fmt.Println("Server running on :8080")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	// ------------------------------------------------
	// OS Signals
	// ------------------------------------------------

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// Wait for CTRL+C.
	<-signalChan

	fmt.Println("\nShutdown signal received")

	// ------------------------------------------------
	// Stop accepting new work
	// ------------------------------------------------

	cancel()

	// ------------------------------------------------
	// Shutdown HTTP server
	// ------------------------------------------------

	shutdownCtx, shutdownCancel :=
		context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("server shutdown error: %v\n", err)
	}

	// ------------------------------------------------
	// Wait for workers
	// ------------------------------------------------

	workerPool.Wait()

	fmt.Println("All workers stopped")
	fmt.Println("Application stopped")
}
