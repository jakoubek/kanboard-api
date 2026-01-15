package kanboard

import "context"

// BoardScope provides fluent project-scoped operations.
type BoardScope struct {
	client    *Client
	projectID int
}

// Board returns a BoardScope for fluent project-scoped operations.
func (c *Client) Board(projectID int) *BoardScope {
	return &BoardScope{
		client:    c,
		projectID: projectID,
	}
}

// GetColumns returns all columns for the project, sorted by position.
func (b *BoardScope) GetColumns(ctx context.Context) ([]Column, error) {
	return b.client.GetColumns(ctx, b.projectID)
}

// GetCategories returns all categories for the project.
func (b *BoardScope) GetCategories(ctx context.Context) ([]Category, error) {
	return b.client.GetAllCategories(ctx, b.projectID)
}

// GetTasks returns all tasks for the project with the specified status.
func (b *BoardScope) GetTasks(ctx context.Context, status TaskStatus) ([]Task, error) {
	return b.client.GetAllTasks(ctx, b.projectID, status)
}

// SearchTasks searches for tasks in the project using Kanboard query syntax.
func (b *BoardScope) SearchTasks(ctx context.Context, query string) ([]Task, error) {
	return b.client.SearchTasks(ctx, b.projectID, query)
}

// CreateTask creates a new task in the project.
// The ProjectID field in the request is overwritten with the board's project ID.
func (b *BoardScope) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	req.ProjectID = b.projectID
	return b.client.CreateTask(ctx, req)
}
