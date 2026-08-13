package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/talhag3/go-job-runner/internal/domain"
	"github.com/talhag3/go-job-runner/internal/service"
)

type JobHandler struct {
	service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{
		service: service,
	}
}

func (h *JobHandler) Create(w http.ResponseWriter, r *http.Request) {
	var job domain.Job

	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"data":    job,
		"message": "job created",
		"error":   nil,
	})
}

func (h *JobHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.GetAll(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"data":    jobs,
		"message": "jobs retrieved",
		"error":   nil,
	})
}

func (h *JobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/jobs/")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "invalid job ID", http.StatusBadRequest)
		return
	}

	job, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"data":    job,
		"message": "job retrieved",
		"error":   nil,
	})
}
