package kanboard

import (
	"context"
	"fmt"
)

// GetAllCategories returns all categories for a project.
func (c *Client) GetAllCategories(ctx context.Context, projectID int) ([]Category, error) {
	params := map[string]int{"project_id": projectID}

	var result []Category
	if err := c.call(ctx, "getAllCategories", params, &result); err != nil {
		return nil, fmt.Errorf("getAllCategories: %w", err)
	}

	return result, nil
}

// GetCategory returns a category by its ID.
// Returns ErrCategoryNotFound if the category does not exist.
func (c *Client) GetCategory(ctx context.Context, categoryID int) (*Category, error) {
	params := map[string]int{"category_id": categoryID}

	var result *Category
	if err := c.call(ctx, "getCategory", params, &result); err != nil {
		return nil, fmt.Errorf("getCategory: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("%w: category %d", ErrCategoryNotFound, categoryID)
	}

	return result, nil
}
