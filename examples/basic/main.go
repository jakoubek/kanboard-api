// Example: basic demonstrates basic client setup and simple operations.
package main

import (
	"context"
	"fmt"
	"log"
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

	// Create a client with API token authentication
	client := kanboard.NewClient(baseURL).
		WithAPIToken(token).
		WithTimeout(30 * time.Second)

	ctx := context.Background()

	// Get all projects
	projects, err := client.GetAllProjects(ctx)
	if err != nil {
		log.Fatalf("Failed to get projects: %v", err)
	}

	fmt.Printf("Found %d projects:\n", len(projects))
	for _, p := range projects {
		fmt.Printf("  - [%d] %s\n", p.ID, p.Name)
	}

	if len(projects) == 0 {
		fmt.Println("No projects found. Create one in Kanboard first.")
		return
	}

	// Use the first project
	projectID := int(projects[0].ID)

	// Get all active tasks for the project
	tasks, err := client.GetAllTasks(ctx, projectID, kanboard.StatusActive)
	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
	}

	fmt.Printf("\nActive tasks in project %d:\n", projectID)
	for _, t := range tasks {
		fmt.Printf("  - [%d] %s\n", t.ID, t.Title)
	}

	// Create a simple task using the direct API
	newTask, err := client.CreateTask(ctx, kanboard.CreateTaskRequest{
		Title:       "Task created via API",
		ProjectID:   projectID,
		Description: "This task was created using the kanboard-api library.",
	})
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("\nCreated task: [%d] %s\n", newTask.ID, newTask.Title)

	// Get task details
	task, err := client.GetTask(ctx, int(newTask.ID))
	if err != nil {
		log.Fatalf("Failed to get task: %v", err)
	}

	fmt.Printf("Task details:\n")
	fmt.Printf("  Title: %s\n", task.Title)
	fmt.Printf("  Description: %s\n", task.Description)
	fmt.Printf("  Column ID: %d\n", task.ColumnID)
	fmt.Printf("  Created: %s\n", task.DateCreation.Time.Format(time.RFC3339))
}
