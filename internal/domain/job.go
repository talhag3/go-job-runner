package domain

type Job struct {
	ID      int64
	Type    string
	Payload map[string]any
	Status  string
}

func (j *Job) MarkProcessing() {
	j.Status = "processing"
}
func (j *Job) MarkCompleted() {
	j.Status = "completed"
}
func (j *Job) MarkFailed() {
	j.Status = "failed"
}
