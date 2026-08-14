# 🚀 Go Concurrent Job Runner

A clean, production-style Go application built to learn and master core Go programming concepts, concurrency patterns, synchronization, and software architecture.

---

## 🎯 Purpose of This Project

I built this project to gain hands-on experience with Go's fundamental primitives and concurrency model by implementing an in-memory asynchronous job execution engine with an HTTP REST API.

Through this project, I explored:
- How to structure a Go backend using **Clean Layered Architecture** (Domain, Repository, Service, Worker Pool, Handler).
- Managing concurrent background jobs safely with **Goroutines** and **Channels**.
- Thread-safe memory management using **RWMutex locks**.
- Controlling task execution and graceful application shutdown using **Context** and **WaitGroup**.

---

## 💡 Key Go Concepts Covered

### 1. 🧵 Goroutines (`go func()`)
Used for non-blocking asynchronous execution, running the HTTP server independently and powering worker loops in the background.

### 2. 🔀 Channels & Select (`chan`, `select`)
- **Buffered Channels:** Acts as an in-memory job queue (`chan *domain.Job`) that allows job submissions without immediately blocking caller routines.
- **`select` Statement:** Multiplexes channel operations to handle incoming jobs alongside cancellation signals (`ctx.Done()`).

### 3. ⏳ Concurrency Synchronization
- **`sync.WaitGroup`:** Ensures all active worker goroutines complete their work safely before the application terminates.
- **`sync.RWMutex`:** Protects the shared in-memory `map[int64]*domain.Job` against data races during concurrent read/write operations.

### 4. 📍 Pointers & Struct Receivers
Used `*domain.Job` pointers to pass memory references, enabling mutations of job status (`MarkProcessing`, `MarkCompleted`, `MarkFailed`) across repository and worker components.

### 5. 🌐 Context & Graceful Shutdown (`context.Context`)
Integrates OS signal trapping (`SIGINT`, `SIGTERM`) to trigger context cancellation, cleanly stopping workers and closing the HTTP server without losing in-flight jobs.

---

## 🏗️ Architecture

```text
HTTP Request (POST/GET)
        │
        ▼
   [ JobHandler ]
        │
        ▼
   [ JobService ]
   ├── Stores Job ────► [ MemoryJobRepository ] (sync.RWMutex + map)
   └── Submits Job ───► [ WorkerPool ] (Buffered Channel)
                              │
                    ┌─────────┼─────────┐
                    ▼         ▼         ▼
                Worker 1   Worker 2   Worker 3  (Goroutines)
                    │         │         │
                    └─────────┴─────────┘
                              │
                              ▼
                     [ DefaultExecutor ]
```

---

## 🚀 Getting Started

### Prerequisites
- [Go 1.22+](https://go.dev/doc/install) installed.

### Run the Application
```bash
go run cmd/server/main.go
```

### API Examples

#### 1. Create a Job
```bash
curl -X POST http://localhost:8080/jobs \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "type": "EMAIL_NOTIFICATION", "payload": {"to": "user@example.com"}}'
```

#### 2. Get All Jobs
```bash
curl http://localhost:8080/jobs
```

#### 3. Get Job by ID
```bash
curl http://localhost:8080/jobs/1
```

#### 4. Test Graceful Shutdown
Press `CTRL+C` in the terminal while jobs are running to see workers complete gracefully before exit.