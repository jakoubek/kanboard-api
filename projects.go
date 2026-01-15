package kanboard

import (
	"context"
	"fmt"
)

// GetAllProjects returns all projects accessible to the authenticated user.
func (c *Client) GetAllProjects(ctx context.Context) ([]Project, error) {
	var result []Project
	if err := c.call(ctx, "getAllProjects", nil, &result); err != nil {
		return nil, fmt.Errorf("getAllProjects: %w", err)
	}
	return result, nil
}

// GetProjectByID returns a project by its ID.
// Returns ErrProjectNotFound if the project does not exist.
func (c *Client) GetProjectByID(ctx context.Context, projectID int) (*Project, error) {
	params := map[string]int{"project_id": projectID}

	var result *Project
	if err := c.call(ctx, "getProjectById", params, &result); err != nil {
		return nil, fmt.Errorf("getProjectById: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: project %d", ErrProjectNotFound, projectID)
	}

	return result, nil
}

// GetProjectByName returns a project by its name.
// Returns ErrProjectNotFound if the project does not exist.
func (c *Client) GetProjectByName(ctx context.Context, name string) (*Project, error) {
	params := map[string]string{"name": name}

	var result *Project
	if err := c.call(ctx, "getProjectByName", params, &result); err != nil {
		return nil, fmt.Errorf("getProjectByName: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: project %q", ErrProjectNotFound, name)
	}

	return result, nil
}
