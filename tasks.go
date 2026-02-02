package kanboard

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
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

// GetTaskByReference returns a task by its external reference within a project.
// Returns ErrTaskNotFound if no task matches the reference.
func (c *Client) GetTaskByReference(ctx context.Context, projectID int, reference string) (*Task, error) {
	params := map[string]any{"project_id": projectID, "reference": reference}

	var result *Task
	if err := c.call(ctx, "getTaskByReference", params, &result); err != nil {
		return nil, fmt.Errorf("getTaskByReference: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: reference %q in project %d", ErrTaskNotFound, reference, projectID)
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
	var taskID IntOrFalse
	if err := c.call(ctx, "createTask", req, &taskID); err != nil {
		return nil, fmt.Errorf("createTask: %w", err)
	}

	if taskID == 0 {
		return nil, fmt.Errorf("createTask: failed to create task")
	}

	// Fetch the created task to return full details
	return c.GetTask(ctx, int(taskID))
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

// UpdateTaskFromParams updates an existing task using TaskUpdateParams.
// This provides a fluent interface for task updates.
func (c *Client) UpdateTaskFromParams(ctx context.Context, taskID int, params *TaskUpdateParams) error {
	req := params.toUpdateTaskRequest(taskID)
	return c.UpdateTask(ctx, req)
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

// SearchTasks searches for tasks in a project using Kanboard's query syntax.
// The query supports filters like: status:open, assignee:me, color:red, etc.
func (c *Client) SearchTasks(ctx context.Context, projectID int, query string) ([]Task, error) {
	params := map[string]any{
		"project_id": projectID,
		"query":      query,
	}

	var result []Task
	if err := c.call(ctx, "searchTasks", params, &result); err != nil {
		return nil, fmt.Errorf("searchTasks: %w", err)
	}

	return result, nil
}

// MoveTaskPosition moves a task to a specific position within a column and swimlane.
// Use position=1 for first position, position=0 to append at end.
func (c *Client) MoveTaskPosition(ctx context.Context, projectID, taskID, columnID, position, swimlaneID int) error {
	params := map[string]int{
		"project_id":  projectID,
		"task_id":     taskID,
		"column_id":   columnID,
		"position":    position,
		"swimlane_id": swimlaneID,
	}

	var success bool
	if err := c.call(ctx, "moveTaskPosition", params, &success); err != nil {
		return fmt.Errorf("moveTaskPosition: %w", err)
	}

	if !success {
		return &OperationFailedError{
			Operation: fmt.Sprintf("moveTaskPosition(task=%d, column=%d, project=%d)", taskID, columnID, projectID),
			Hints: []string{
				"task may not exist",
				"column may not belong to project",
				"insufficient permissions",
				"task may already be in target position",
			},
		}
	}

	return nil
}

// MoveTaskToProject moves a task to a different project.
func (c *Client) MoveTaskToProject(ctx context.Context, taskID, projectID int) error {
	params := map[string]int{
		"task_id":    taskID,
		"project_id": projectID,
	}

	var success bool
	if err := c.call(ctx, "moveTaskToProject", params, &success); err != nil {
		return fmt.Errorf("moveTaskToProject: %w", err)
	}

	if !success {
		return &OperationFailedError{
			Operation: fmt.Sprintf("moveTaskToProject(task=%d, project=%d)", taskID, projectID),
			Hints: []string{
				"task may not exist",
				"target project may not exist",
				"insufficient permissions",
			},
		}
	}

	return nil
}

// SearchTasksGlobally searches for tasks across all accessible projects.
// The search is executed in parallel across all projects using errgroup.
// If any project search fails, all ongoing searches are cancelled.
func (c *Client) SearchTasksGlobally(ctx context.Context, query string) ([]Task, error) {
	// Get all accessible projects
	projects, err := c.GetAllProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("searchTasksGlobally: %w", err)
	}

	if len(projects) == 0 {
		return []Task{}, nil
	}

	// Use errgroup for parallel execution with context cancellation
	g, ctx := errgroup.WithContext(ctx)

	// Slice to store results from each project (one per project, thread-safe by index)
	results := make([][]Task, len(projects))

	// Launch parallel searches
	for i, project := range projects {
		i, projectID := i, int(project.ID)
		g.Go(func() error {
			tasks, err := c.SearchTasks(ctx, projectID, query)
			if err != nil {
				return err
			}
			results[i] = tasks
			return nil
		})
	}

	// Wait for all goroutines
	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("searchTasksGlobally: %w", err)
	}

	// Aggregate results
	var allTasks []Task
	for _, tasks := range results {
		allTasks = append(allTasks, tasks...)
	}

	return allTasks, nil
}
