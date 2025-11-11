package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/antlabs/gurl/internal/stats"
)

// TaskStatus represents the status of a benchmark task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Task represents a benchmark task
type Task struct {
	ID          string                 `json:"id"`
	Status      TaskStatus             `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Results     *stats.Results         `json:"results,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// TaskManager manages benchmark tasks
type TaskManager struct {
	mu    sync.RWMutex
	tasks map[string]*Task
}

// NewTaskManager creates a new task manager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*Task),
	}
}

// CreateTask creates a new task
func (tm *TaskManager) CreateTask(id string, config map[string]interface{}) *Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task := &Task{
		ID:        id,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		Config:    config,
	}

	tm.tasks[id] = task
	return task
}

// GetTask retrieves a task by ID
func (tm *TaskManager) GetTask(id string) (*Task, bool) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	task, exists := tm.tasks[id]
	if !exists {
		return nil, false
	}

	// Return a copy to avoid race conditions
	return &Task{
		ID:          task.ID,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		StartedAt:   task.StartedAt,
		CompletedAt: task.CompletedAt,
		Error:       task.Error,
		Results:     task.Results, // Results is already thread-safe
		Config:      task.Config,
	}, true
}

// UpdateTaskStatus updates the status of a task
func (tm *TaskManager) UpdateTaskStatus(id string, status TaskStatus) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[id]
	if !exists {
		return fmt.Errorf("task %s not found", id)
	}

	task.Status = status
	now := time.Now()

	switch status {
	case TaskStatusRunning:
		if task.StartedAt == nil {
			task.StartedAt = &now
		}
	case TaskStatusCompleted, TaskStatusFailed:
		if task.CompletedAt == nil {
			task.CompletedAt = &now
		}
	}

	return nil
}

// SetTaskError sets the error message for a task
func (tm *TaskManager) SetTaskError(id string, err error) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[id]
	if !exists {
		return fmt.Errorf("task %s not found", id)
	}

	task.Error = err.Error()
	task.Status = TaskStatusFailed
	now := time.Now()
	if task.CompletedAt == nil {
		task.CompletedAt = &now
	}

	return nil
}

// SetTaskResults sets the results for a completed task
func (tm *TaskManager) SetTaskResults(id string, results *stats.Results) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task, exists := tm.tasks[id]
	if !exists {
		return fmt.Errorf("task %s not found", id)
	}

	task.Results = results
	task.Status = TaskStatusCompleted
	now := time.Now()
	if task.CompletedAt == nil {
		task.CompletedAt = &now
	}

	return nil
}

// ListTasks returns all tasks
func (tm *TaskManager) ListTasks() []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, &Task{
			ID:          task.ID,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
			StartedAt:   task.StartedAt,
			CompletedAt: task.CompletedAt,
			Error:       task.Error,
			Config:      task.Config,
			// Don't include Results in list to save memory
		})
	}

	return tasks
}

// RunTask runs a benchmark task asynchronously
func (tm *TaskManager) RunTask(ctx context.Context, id string, runner func(context.Context) (*stats.Results, error)) {
	// Update status to running
	tm.UpdateTaskStatus(id, TaskStatusRunning)

	// Run benchmark in goroutine
	go func() {
		results, err := runner(ctx)
		if err != nil {
			tm.SetTaskError(id, err)
			return
		}

		tm.SetTaskResults(id, results)
	}()
}

