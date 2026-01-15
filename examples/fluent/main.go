// Example: fluent demonstrates the fluent API for task management.
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	kanboard "code.beautifulmachines.dev/jakoubek/kanboard-api"
)

func main() {
	// Read configuration from environment
	baseURL := os.Getenv("KANBOARD_URL")
	token := os.Getenv("KANBOARD_TOKEN")

	if baseURL == "" || token == "" {
		log.Fatal("Set KANBOARD_URL and KANBOARD_TOKEN environment variables")
	}

	// Create a fully configured client
	client := kanboard.NewClient(baseURL).
		WithAPIToken(token).
		WithTimeout(60 * time.Second).
		WithLogger(slog.Default()) // Enable debug logging

	ctx := context.Background()

	// Get the first project
	projects, err := client.GetAllProjects(ctx)
	if err != nil {
		log.Fatalf("Failed to get projects: %v", err)
	}
	if len(projects) == 0 {
		log.Fatal("No projects found")
	}
	projectID := int(projects[0].ID)

	// Use BoardScope for project-scoped operations
	board := client.Board(projectID)

	// Get columns to understand the workflow
	columns, err := board.GetColumns(ctx)
	if err != nil {
		log.Fatalf("Failed to get columns: %v", err)
	}
	fmt.Println("Columns in project:")
	for _, c := range columns {
		fmt.Printf("  - [%d] %s (position: %d)\n", c.ID, c.Title, c.Position)
	}

	// Create a task using TaskParams (fluent builder)
	params := kanboard.NewTask("Feature: User Dashboard").
		WithDescription("Implement the user dashboard with activity feed").
		WithPriority(2).
		WithScore(5).
		WithTags("feature", "frontend", "v2.0").
		WithDueDate(time.Now().Add(14 * 24 * time.Hour))

	task, err := board.CreateTaskFromParams(ctx, params)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}
	fmt.Printf("\nCreated task: [%d] %s\n", task.ID, task.Title)

	// Use TaskScope for task-scoped operations
	taskScope := client.Task(int(task.ID))

	// Add more tags
	if err := taskScope.AddTag(ctx, "priority"); err != nil {
		log.Fatalf("Failed to add tag: %v", err)
	}
	fmt.Println("Added 'priority' tag")

	// Get current tags
	tags, err := taskScope.GetTags(ctx)
	if err != nil {
		log.Fatalf("Failed to get tags: %v", err)
	}
	fmt.Printf("Current tags: %v\n", tags)

	// Update the task using TaskUpdateParams
	updates := kanboard.NewTaskUpdate().
		SetTitle("Feature: Enhanced User Dashboard").
		SetDescription("Implement user dashboard with activity feed and notifications").
		SetPriority(1)

	if err := taskScope.Update(ctx, updates); err != nil {
		log.Fatalf("Failed to update task: %v", err)
	}
	fmt.Println("Updated task title and priority")

	// Move task to next column
	if err := taskScope.MoveToNextColumn(ctx); err != nil {
		if err == kanboard.ErrAlreadyInLastColumn {
			fmt.Println("Task is already in the last column")
		} else {
			log.Fatalf("Failed to move task: %v", err)
		}
	} else {
		fmt.Println("Moved task to next column")
	}

	// Get updated task details
	updated, err := taskScope.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}
	fmt.Printf("\nUpdated task:\n")
	fmt.Printf("  Title: %s\n", updated.Title)
	fmt.Printf("  Column: %d\n", updated.ColumnID)
	fmt.Printf("  Priority: %d\n", updated.Priority)
}
