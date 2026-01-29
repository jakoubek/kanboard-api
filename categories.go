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

// CreateCategory creates a new category and returns its ID.
func (c *Client) CreateCategory(ctx context.Context, projectID int, name string, colorID string) (int, error) {
	params := map[string]interface{}{
		"project_id": projectID,
		"name":       name,
	}
	if colorID != "" {
		params["color_id"] = colorID
	}

	var result IntOrFalse
	if err := c.call(ctx, "createCategory", params, &result); err != nil {
		return 0, fmt.Errorf("createCategory: %w", err)
	}

	if int(result) == 0 {
		return 0, fmt.Errorf("createCategory: failed to create category %q", name)
	}

	return int(result), nil
}

// UpdateCategory updates a category's name and optionally its color.
func (c *Client) UpdateCategory(ctx context.Context, categoryID int, name string, colorID string) error {
	params := map[string]interface{}{
		"id":   categoryID,
		"name": name,
	}
	if colorID != "" {
		params["color_id"] = colorID
	}

	var result bool
	if err := c.call(ctx, "updateCategory", params, &result); err != nil {
		return fmt.Errorf("updateCategory: %w", err)
	}

	if !result {
		return fmt.Errorf("updateCategory: failed to update category %d", categoryID)
	}

	return nil
}

// RemoveCategory deletes a category.
func (c *Client) RemoveCategory(ctx context.Context, categoryID int) error {
	params := map[string]int{"category_id": categoryID}

	var result bool
	if err := c.call(ctx, "removeCategory", params, &result); err != nil {
		return fmt.Errorf("removeCategory: %w", err)
	}

	if !result {
		return fmt.Errorf("removeCategory: failed to remove category %d", categoryID)
	}

	return nil
}

// GetCategoryByName returns a category by name within a project.
// Returns ErrCategoryNotFound if no category matches.
func (c *Client) GetCategoryByName(ctx context.Context, projectID int, name string) (*Category, error) {
	categories, err := c.GetAllCategories(ctx, projectID)
	if err != nil {
		return nil, err
	}

	for i := range categories {
		if categories[i].Name == name {
			return &categories[i], nil
		}
	}

	return nil, fmt.Errorf("%w: category %q in project %d", ErrCategoryNotFound, name, projectID)
}
