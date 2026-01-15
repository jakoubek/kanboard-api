package kanboard

import (
	"context"
	"fmt"
)

// GetAllComments returns all comments for a task.
func (c *Client) GetAllComments(ctx context.Context, taskID int) ([]Comment, error) {
	params := map[string]int{"task_id": taskID}

	var result []Comment
	if err := c.call(ctx, "getAllComments", params, &result); err != nil {
		return nil, fmt.Errorf("getAllComments: %w", err)
	}

	return result, nil
}

// GetComment returns a comment by its ID.
// Returns ErrCommentNotFound if the comment does not exist.
func (c *Client) GetComment(ctx context.Context, commentID int) (*Comment, error) {
	params := map[string]int{"comment_id": commentID}

	var result *Comment
	if err := c.call(ctx, "getComment", params, &result); err != nil {
		return nil, fmt.Errorf("getComment: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: comment %d", ErrCommentNotFound, commentID)
	}

	return result, nil
}

// CreateComment creates a new comment on a task and returns the created comment.
func (c *Client) CreateComment(ctx context.Context, taskID, userID int, content string) (*Comment, error) {
	params := map[string]any{
		"task_id": taskID,
		"user_id": userID,
		"content": content,
	}

	var commentID int
	if err := c.call(ctx, "createComment", params, &commentID); err != nil {
		return nil, fmt.Errorf("createComment: %w", err)
	}

	if commentID == 0 {
		return nil, fmt.Errorf("createComment: failed to create comment")
	}

	// Fetch the created comment to return full details
	return c.GetComment(ctx, commentID)
}

// UpdateComment updates the content of a comment.
func (c *Client) UpdateComment(ctx context.Context, commentID int, content string) error {
	params := map[string]any{
		"id":      commentID,
		"content": content,
	}

	var success bool
	if err := c.call(ctx, "updateComment", params, &success); err != nil {
		return fmt.Errorf("updateComment: %w", err)
	}

	if !success {
		return fmt.Errorf("updateComment: update failed")
	}

	return nil
}

// RemoveComment deletes a comment.
func (c *Client) RemoveComment(ctx context.Context, commentID int) error {
	params := map[string]int{"comment_id": commentID}

	var success bool
	if err := c.call(ctx, "removeComment", params, &success); err != nil {
		return fmt.Errorf("removeComment: %w", err)
	}

	if !success {
		return fmt.Errorf("removeComment: delete failed")
	}

	return nil
}
