package kanboard

import (
	"context"
	"fmt"
)

// GetAllLinks returns all link types configured in the Kanboard instance
// (e.g. "relates to", "blocks", "is blocked by"). The IDs are instance-specific:
// they match Kanboard's default seed only as long as no link type has been
// added, removed or reordered.
func (c *Client) GetAllLinks(ctx context.Context) ([]Link, error) {
	var result []Link
	if err := c.call(ctx, "getAllLinks", nil, &result); err != nil {
		return nil, fmt.Errorf("getAllLinks: %w", err)
	}

	return result, nil
}

// GetOppositeLinkID returns the ID of the link type pointing in the opposite
// direction, or 0 if the link type has no opposite. GetAllLinks already carries
// this information in Link.OppositeID, so this is only needed for a lookup of a
// single link type.
func (c *Client) GetOppositeLinkID(ctx context.Context, linkID int) (int, error) {
	params := map[string]int{"link_id": linkID}

	var result StringInt
	if err := c.call(ctx, "getOppositeLinkId", params, &result); err != nil {
		return 0, fmt.Errorf("getOppositeLinkId: %w", err)
	}

	return int(result), nil
}

// GetAllTaskLinks returns all links for a task.
func (c *Client) GetAllTaskLinks(ctx context.Context, taskID int) ([]TaskLink, error) {
	params := map[string]int{"task_id": taskID}

	var result []TaskLink
	if err := c.call(ctx, "getAllTaskLinks", params, &result); err != nil {
		return nil, fmt.Errorf("getAllTaskLinks: %w", err)
	}

	return result, nil
}

// CreateTaskLink creates a link between two tasks.
// The linkID specifies the type of relationship (e.g., "blocks", "is blocked by").
// Returns the ID of the created link.
func (c *Client) CreateTaskLink(ctx context.Context, taskID, oppositeTaskID, linkID int) (int, error) {
	params := map[string]int{
		"task_id":          taskID,
		"opposite_task_id": oppositeTaskID,
		"link_id":          linkID,
	}

	var result IntOrFalse
	if err := c.call(ctx, "createTaskLink", params, &result); err != nil {
		return 0, fmt.Errorf("createTaskLink: %w", err)
	}

	if result == 0 {
		return 0, fmt.Errorf("createTaskLink: failed to create link")
	}

	return int(result), nil
}

// RemoveTaskLink deletes a task link.
func (c *Client) RemoveTaskLink(ctx context.Context, taskLinkID int) error {
	params := map[string]int{"task_link_id": taskLinkID}

	var success bool
	if err := c.call(ctx, "removeTaskLink", params, &success); err != nil {
		return fmt.Errorf("removeTaskLink: %w", err)
	}

	if !success {
		return fmt.Errorf("removeTaskLink: delete failed")
	}

	return nil
}
