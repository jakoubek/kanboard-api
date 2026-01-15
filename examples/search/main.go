// Example: search demonstrates search functionality.
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

	client := kanboard.NewClient(baseURL).
		WithAPIToken(token).
		WithTimeout(30 * time.Second)

	ctx := context.Background()

	// Get the first project for project-specific search
	projects, err := client.GetAllProjects(ctx)
	if err != nil {
		log.Fatalf("Failed to get projects: %v", err)
	}
	if len(projects) == 0 {
		log.Fatal("No projects found")
	}
	projectID := int(projects[0].ID)
	fmt.Printf("Using project: %s (ID: %d)\n\n", projects[0].Name, projectID)

	// Project-specific search using Kanboard query syntax
	// See: https://docs.kanboard.org/v1/user/search/
	fmt.Println("=== Project-Specific Search ===")

	// Search for open tasks
	tasks, err := client.SearchTasks(ctx, projectID, "status:open")
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}
	fmt.Printf("Open tasks: %d\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("  - [%d] %s\n", t.ID, t.Title)
	}

	// Search by title keyword
	tasks, err = client.Board(projectID).SearchTasks(ctx, "title:feature")
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}
	fmt.Printf("\nTasks with 'feature' in title: %d\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("  - [%d] %s\n", t.ID, t.Title)
	}

	// Search overdue tasks
	tasks, err = client.SearchTasks(ctx, projectID, "due:<today")
	if err != nil {
		log.Fatalf("Failed to search: %v", err)
	}
	fmt.Printf("\nOverdue tasks: %d\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("  - [%d] %s (due: %s)\n", t.ID, t.Title, t.DateDue.Time.Format("2006-01-02"))
	}

	// Global search across ALL projects
	fmt.Println("\n=== Global Search (All Projects) ===")

	// This searches in parallel across all accessible projects
	allTasks, err := client.SearchTasksGlobally(ctx, "status:open")
	if err != nil {
		log.Fatalf("Failed to global search: %v", err)
	}
	fmt.Printf("Open tasks across all projects: %d\n", len(allTasks))

	// Group results by project
	byProject := make(map[int][]kanboard.Task)
	for _, t := range allTasks {
		pid := int(t.ProjectID)
		byProject[pid] = append(byProject[pid], t)
	}

	for pid, tasks := range byProject {
		fmt.Printf("\n  Project %d: %d tasks\n", pid, len(tasks))
		for _, t := range tasks {
			fmt.Printf("    - [%d] %s\n", t.ID, t.Title)
		}
	}

	// Search examples with Kanboard query syntax
	fmt.Println("\n=== Query Syntax Examples ===")
	fmt.Println("Kanboard supports rich query syntax:")
	fmt.Println("  status:open              - Open tasks")
	fmt.Println("  status:closed            - Closed tasks")
	fmt.Println("  assignee:me              - Tasks assigned to current user")
	fmt.Println("  due:today                - Tasks due today")
	fmt.Println("  due:<tomorrow            - Tasks due before tomorrow")
	fmt.Println("  due:>2024-01-01          - Tasks due after date")
	fmt.Println("  title:\"bug fix\"          - Tasks with exact title match")
	fmt.Println("  color:red                - Tasks with red color")
	fmt.Println("  category:\"Bug\"           - Tasks in category")
	fmt.Println("  tag:urgent               - Tasks with tag")
	fmt.Println("  priority:3               - Tasks with priority")
	fmt.Println()
	fmt.Println("Combine queries with spaces (AND) or | (OR):")
	fmt.Println("  status:open assignee:me  - Open tasks assigned to me")
	fmt.Println("  status:open | status:closed - All tasks")
}
