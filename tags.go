package kanboard

import (
	"context"
	"fmt"
)

// GetTaskTags returns the tags assigned to a task as a map of tagID to tag name.
func (c *Client) GetTaskTags(ctx context.Context, taskID int) (map[int]string, error) {
	params := map[string]int{"task_id": taskID}

	// Kanboard returns map[string]string where keys are string tag IDs
	var result map[string]string
	if err := c.call(ctx, "getTaskTags", params, &result); err != nil {
		return nil, fmt.Errorf("getTaskTags: %w", err)
	}

	// Convert string keys to int
	tags := make(map[int]string, len(result))
	for idStr, name := range result {
		var id int
		if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
			continue // Skip invalid IDs
		}
		tags[id] = name
	}

	return tags, nil
}

// SetTaskTags sets the tags for a task, replacing all existing tags.
// Tags are specified by name. Non-existent tags will be auto-created.
func (c *Client) SetTaskTags(ctx context.Context, projectID, taskID int, tags []string) error {
	params := map[string]interface{}{
		"project_id": projectID,
		"task_id":    taskID,
		"tags":       tags,
	}

	var result bool
	if err := c.call(ctx, "setTaskTags", params, &result); err != nil {
		return fmt.Errorf("setTaskTags: %w", err)
	}

	return nil
}

// GetAllTags returns all tags in the system.
func (c *Client) GetAllTags(ctx context.Context) ([]Tag, error) {
	var result []Tag
	if err := c.call(ctx, "getAllTags", nil, &result); err != nil {
		return nil, fmt.Errorf("getAllTags: %w", err)
	}
	return result, nil
}

// GetTagsByProject returns all tags for a specific project.
func (c *Client) GetTagsByProject(ctx context.Context, projectID int) ([]Tag, error) {
	params := map[string]int{"project_id": projectID}

	var result []Tag
	if err := c.call(ctx, "getTagsByProject", params, &result); err != nil {
		return nil, fmt.Errorf("getTagsByProject: %w", err)
	}
	return result, nil
}

// CreateTag creates a new tag in a project and returns the tag ID.
func (c *Client) CreateTag(ctx context.Context, projectID int, name, colorID string) (int, error) {
	params := map[string]interface{}{
		"project_id": projectID,
		"tag":        name,
	}
	if colorID != "" {
		params["color_id"] = colorID
	}

	var result IntOrFalse
	if err := c.call(ctx, "createTag", params, &result); err != nil {
		return 0, fmt.Errorf("createTag: %w", err)
	}
	return int(result), nil
}

// UpdateTag updates an existing tag's name and/or color.
func (c *Client) UpdateTag(ctx context.Context, tagID int, name, colorID string) error {
	params := map[string]interface{}{
		"tag_id": tagID,
		"tag":    name,
	}
	if colorID != "" {
		params["color_id"] = colorID
	}

	var result bool
	if err := c.call(ctx, "updateTag", params, &result); err != nil {
		return fmt.Errorf("updateTag: %w", err)
	}
	return nil
}

// RemoveTag deletes a tag from the system.
func (c *Client) RemoveTag(ctx context.Context, tagID int) error {
	params := map[string]int{"tag_id": tagID}

	var result bool
	if err := c.call(ctx, "removeTag", params, &result); err != nil {
		return fmt.Errorf("removeTag: %w", err)
	}
	return nil
}
