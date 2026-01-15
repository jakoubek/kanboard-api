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

// GetTags returns the tags assigned to this task as a map of tagID to tag name.
func (t *TaskScope) GetTags(ctx context.Context) (map[int]string, error) {
	return t.client.GetTaskTags(ctx, t.taskID)
}

// SetTags sets the tags for this task, replacing all existing tags.
// Tags are specified by name. Non-existent tags will be auto-created.
//
// WARNING: This operation is not atomic. Concurrent tag modifications may cause data loss.
func (t *TaskScope) SetTags(ctx context.Context, tags ...string) error {
	task, err := t.Get(ctx)
	if err != nil {
		return err
	}
	return t.client.SetTaskTags(ctx, int(task.ProjectID), t.taskID, tags)
}

// ClearTags removes all tags from this task.
//
// WARNING: This operation is not atomic. Concurrent tag modifications may cause data loss.
func (t *TaskScope) ClearTags(ctx context.Context) error {
	return t.SetTags(ctx)
}

// AddTag adds a tag to this task using a read-modify-write pattern.
// If the tag already exists on the task, this is a no-op (idempotent).
//
// WARNING: This operation is not atomic. Concurrent tag modifications may cause data loss.
func (t *TaskScope) AddTag(ctx context.Context, tag string) error {
	// Get task info for project_id
	task, err := t.Get(ctx)
	if err != nil {
		return err
	}

	// Get current tags
	currentTags, err := t.client.GetTaskTags(ctx, t.taskID)
	if err != nil {
		return err
	}

	// Check if tag already exists (idempotent)
	for _, existingTag := range currentTags {
		if existingTag == tag {
			return nil
		}
	}

	// Build new tag list
	tagNames := make([]string, 0, len(currentTags)+1)
	for _, name := range currentTags {
		tagNames = append(tagNames, name)
	}
	tagNames = append(tagNames, tag)

	// Set updated tags
	return t.client.SetTaskTags(ctx, int(task.ProjectID), t.taskID, tagNames)
}

// RemoveTag removes a tag from this task using a read-modify-write pattern.
// If the tag doesn't exist on the task, this is a no-op (idempotent).
//
// WARNING: This operation is not atomic. Concurrent tag modifications may cause data loss.
func (t *TaskScope) RemoveTag(ctx context.Context, tag string) error {
	// Get task info for project_id
	task, err := t.Get(ctx)
	if err != nil {
		return err
	}

	// Get current tags
	currentTags, err := t.client.GetTaskTags(ctx, t.taskID)
	if err != nil {
		return err
	}

	// Filter out the tag to remove
	tagNames := make([]string, 0, len(currentTags))
	found := false
	for _, name := range currentTags {
		if name == tag {
			found = true
			continue
		}
		tagNames = append(tagNames, name)
	}

	// If tag wasn't found, nothing to do (idempotent)
	if !found {
		return nil
	}

	// Set updated tags
	return t.client.SetTaskTags(ctx, int(task.ProjectID), t.taskID, tagNames)
}

// HasTag checks if this task has a specific tag.
func (t *TaskScope) HasTag(ctx context.Context, tag string) (bool, error) {
	tags, err := t.client.GetTaskTags(ctx, t.taskID)
	if err != nil {
		return false, err
	}

	for _, name := range tags {
		if name == tag {
			return true, nil
		}
	}

	return false, nil
}

// GetComments returns all comments for this task.
func (t *TaskScope) GetComments(ctx context.Context) ([]Comment, error) {
	return t.client.GetAllComments(ctx, t.taskID)
}

// AddComment adds a comment to this task and returns the created comment.
// The userID is the ID of the user creating the comment.
func (t *TaskScope) AddComment(ctx context.Context, userID int, content string) (*Comment, error) {
	return t.client.CreateComment(ctx, t.taskID, userID, content)
}
