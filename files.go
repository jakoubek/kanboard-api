package kanboard

import (
	"context"
	"encoding/base64"
	"fmt"
)

// GetAllTaskFiles returns all files attached to a task.
func (c *Client) GetAllTaskFiles(ctx context.Context, taskID int) ([]TaskFile, error) {
	params := map[string]int{"task_id": taskID}

	var result []TaskFile
	if err := c.call(ctx, "getAllTaskFiles", params, &result); err != nil {
		return nil, fmt.Errorf("getAllTaskFiles: %w", err)
	}

	return result, nil
}

// CreateTaskFile uploads a file to a task.
// The file content is automatically base64 encoded.
// Returns the ID of the created file.
func (c *Client) CreateTaskFile(ctx context.Context, projectID, taskID int, filename string, content []byte) (int, error) {
	params := map[string]any{
		"project_id": projectID,
		"task_id":    taskID,
		"filename":   filename,
		"blob":       base64.StdEncoding.EncodeToString(content),
	}

	var result IntOrFalse
	if err := c.call(ctx, "createTaskFile", params, &result); err != nil {
		return 0, fmt.Errorf("createTaskFile: %w", err)
	}

	if result == 0 {
		return 0, fmt.Errorf("createTaskFile: failed to upload file")
	}

	return int(result), nil
}

// DownloadTaskFile downloads a file's content by its ID.
// The content is returned as raw bytes (decoded from base64).
func (c *Client) DownloadTaskFile(ctx context.Context, fileID int) ([]byte, error) {
	params := map[string]int{"file_id": fileID}

	var result string
	if err := c.call(ctx, "downloadTaskFile", params, &result); err != nil {
		return nil, fmt.Errorf("downloadTaskFile: %w", err)
	}

	// Decode base64 content
	content, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		return nil, fmt.Errorf("downloadTaskFile: failed to decode content: %w", err)
	}

	return content, nil
}

// RemoveTaskFile deletes a file from a task.
func (c *Client) RemoveTaskFile(ctx context.Context, fileID int) error {
	params := map[string]int{"file_id": fileID}

	var success bool
	if err := c.call(ctx, "removeTaskFile", params, &success); err != nil {
		return fmt.Errorf("removeTaskFile: %w", err)
	}

	if !success {
		return fmt.Errorf("removeTaskFile: delete failed")
	}

	return nil
}
