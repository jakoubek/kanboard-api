package kanboard

import (
	"context"
	"fmt"
)

// GetColorList returns the available task colors.
// Returns a map of color_id to display name (e.g., "yellow" -> "Yellow").
func (c *Client) GetColorList(ctx context.Context) (map[string]string, error) {
	var result map[string]string
	if err := c.call(ctx, "getColorList", nil, &result); err != nil {
		return nil, fmt.Errorf("getColorList: %w", err)
	}

	return result, nil
}
