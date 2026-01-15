package kanboard

import (
	"context"
	"fmt"
)

// GetTask returns a task by its ID.
// Returns ErrTaskNotFound if the task does not exist.
func (c *Client) GetTask(ctx context.Context, taskID int) (*Task, error) {
	params := map[string]int{"task_id": taskID}

	var result *Task
	if err := c.call(ctx, "getTask", params, &result); err != nil {
		return nil, fmt.Errorf("getTask: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: task %d", ErrTaskNotFound, taskID)
	}

	return result, nil
}

// GetAllTasks returns all tasks for a project with the specified status.
func (c *Client) GetAllTasks(ctx context.Context, projectID int, status TaskStatus) ([]Task, error) {
	params := map[string]int{
		"project_id": projectID,
		"status_id":  int(status),
	}

	var result []Task
	if err := c.call(ctx, "getAllTasks", params, &result); err != nil {
		return nil, fmt.Errorf("getAllTasks: %w", err)
	}

	return result, nil
}

// CreateTask creates a new task and returns the created task.
func (c *Client) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	var taskID int
	if err := c.call(ctx, "createTask", req, &taskID); err != nil {
		return nil, fmt.Errorf("createTask: %w", err)
	}

	if taskID == 0 {
		return nil, fmt.Errorf("createTask: failed to create task")
	}

	// Fetch the created task to return full details
	return c.GetTask(ctx, taskID)
}

// UpdateTask updates an existing task.
// Only non-nil fields in the request will be updated.
func (c *Client) UpdateTask(ctx context.Context, req UpdateTaskRequest) error {
	var success bool
	if err := c.call(ctx, "updateTask", req, &success); err != nil {
		return fmt.Errorf("updateTask: %w", err)
	}

	if !success {
		return fmt.Errorf("updateTask: update failed")
	}

	return nil
}

// CloseTask closes a task (sets it to inactive).
// Returns ErrTaskClosed if the task is already closed.
func (c *Client) CloseTask(ctx context.Context, taskID int) error {
	params := map[string]int{"task_id": taskID}

	var success bool
	if err := c.call(ctx, "closeTask", params, &success); err != nil {
		return fmt.Errorf("closeTask: %w", err)
	}

	if !success {
		return fmt.Errorf("%w: task %d", ErrTaskClosed, taskID)
	}

	return nil
}

// OpenTask opens a task (sets it to active).
// Returns ErrTaskOpen if the task is already open.
func (c *Client) OpenTask(ctx context.Context, taskID int) error {
	params := map[string]int{"task_id": taskID}

	var success bool
	if err := c.call(ctx, "openTask", params, &success); err != nil {
		return fmt.Errorf("openTask: %w", err)
	}

	if !success {
		return fmt.Errorf("%w: task %d", ErrTaskOpen, taskID)
	}

	return nil
}
