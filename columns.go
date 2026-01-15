package kanboard

import (
	"context"
	"fmt"
	"sort"
)

// GetColumns returns all columns for a project, sorted by position.
func (c *Client) GetColumns(ctx context.Context, projectID int) ([]Column, error) {
	params := map[string]int{"project_id": projectID}

	var result []Column
	if err := c.call(ctx, "getColumns", params, &result); err != nil {
		return nil, fmt.Errorf("getColumns: %w", err)
	}

	// Sort by position
	sort.Slice(result, func(i, j int) bool {
		return int(result[i].Position) < int(result[j].Position)
	})

	return result, nil
}

// GetColumn returns a column by its ID.
// Returns ErrColumnNotFound if the column does not exist.
func (c *Client) GetColumn(ctx context.Context, columnID int) (*Column, error) {
	params := map[string]int{"column_id": columnID}

	var result *Column
	if err := c.call(ctx, "getColumn", params, &result); err != nil {
		return nil, fmt.Errorf("getColumn: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: column %d", ErrColumnNotFound, columnID)
	}

	return result, nil
}
