package kanboard

import "context"

// TaskScope provides fluent task-scoped operations.
type TaskScope struct {
	client *Client
	taskID int
}

// Task returns a TaskScope for fluent task-scoped operations.
func (c *Client) Task(taskID int) *TaskScope {
	return &TaskScope{
		client: c,
		taskID: taskID,
	}
}

// Get returns the task.
func (t *TaskScope) Get(ctx context.Context) (*Task, error) {
	return t.client.GetTask(ctx, t.taskID)
}

// Close closes the task (sets it to inactive).
func (t *TaskScope) Close(ctx context.Context) error {
	return t.client.CloseTask(ctx, t.taskID)
}

// Open opens the task (sets it to active).
func (t *TaskScope) Open(ctx context.Context) error {
	return t.client.OpenTask(ctx, t.taskID)
}

// MoveToColumn moves the task to a different column.
// The task is placed at the end of the column (position=0).
// Requires the project ID to be fetched from the task.
func (t *TaskScope) MoveToColumn(ctx context.Context, columnID int) error {
	task, err := t.Get(ctx)
	if err != nil {
		return err
	}
	return t.client.MoveTaskPosition(ctx, int(task.ProjectID), t.taskID, columnID, 0, int(task.SwimlaneID))
}

// MoveToProject moves the task to a different project.
func (t *TaskScope) MoveToProject(ctx context.Context, projectID int) error {
	return t.client.MoveTaskToProject(ctx, t.taskID, projectID)
}

// Update updates the task using TaskUpdateParams.
func (t *TaskScope) Update(ctx context.Context, params *TaskUpdateParams) error {
	return t.client.UpdateTaskFromParams(ctx, t.taskID, params)
}
