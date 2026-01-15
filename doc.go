// Package kanboard provides a Go client for the Kanboard JSON-RPC API.
//
// This library offers both a fluent API for convenient, chainable operations
// and direct API methods for lower-level access. The client is safe for
// concurrent use by multiple goroutines.
//
// # Quick Start
//
// Create a client with API token authentication:
//
//	client := kanboard.NewClient("https://kanboard.example.com").
//	    WithAPIToken("your-api-token")
//
// # Authentication
//
// The library supports two authentication methods:
//
//   - API Token: Use [Client.WithAPIToken] for token-based auth (recommended)
//   - Basic Auth: Use [Client.WithBasicAuth] for username/password auth
//
// # Fluent API
//
// Use [Client.Board] for project-scoped operations:
//
//	board := client.Board(projectID)
//	tasks, _ := board.GetTasks(ctx, kanboard.StatusActive)
//	task, _ := board.CreateTaskFromParams(ctx,
//	    kanboard.NewTask("Title").WithDescription("Details"))
//
// Use [Client.Task] for task-scoped operations:
//
//	task := client.Task(taskID)
//	task.MoveToNextColumn(ctx)
//	task.AddTag(ctx, "reviewed")
//	task.AddComment(ctx, userID, "Comment text")
//
// # Task Creation
//
// Use [TaskParams] for fluent task creation:
//
//	params := kanboard.NewTask("Title").
//	    WithDescription("Details").
//	    WithPriority(2).
//	    WithTags("urgent", "backend")
//
// # Error Handling
//
// The library provides typed errors and helper functions:
//
//	task, err := client.GetTask(ctx, taskID)
//	if kanboard.IsNotFound(err) {
//	    // Handle not found
//	}
//	if kanboard.IsUnauthorized(err) {
//	    // Handle auth failure
//	}
//
// # Thread Safety
//
// The [Client] is safe for concurrent use. Request IDs are generated atomically.
//
// # Tag Operations
//
// Warning: Kanboard's setTaskTags API replaces all tags. The [TaskScope.AddTag]
// and [TaskScope.RemoveTag] methods use read-modify-write internally and are
// not atomic. Concurrent tag modifications may cause data loss.
package kanboard
