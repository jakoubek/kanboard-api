# kanboard-api

[![Mirror on GitHub](https://img.shields.io/badge/mirror-GitHub-blue)](https://github.com/jakoubek/kanboard-api)
[![Go Reference](https://pkg.go.dev/badge/code.beautifulmachines.dev/jakoubek/kanboard-api.svg)](https://pkg.go.dev/code.beautifulmachines.dev/jakoubek/kanboard-api)
[![Go Report Card](https://goreportcard.com/badge/code.beautifulmachines.dev/jakoubek/kanboard-api)](https://goreportcard.com/report/code.beautifulmachines.dev/jakoubek/kanboard-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Go library for the [Kanboard](https://kanboard.org/) JSON-RPC API. Provides a fluent, chainable API for integrating Kanboard into Go applications.

## Installation

```bash
go get code.beautifulmachines.dev/jakoubek/kanboard-api
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    kanboard "code.beautifulmachines.dev/jakoubek/kanboard-api"
)

func main() {
    ctx := context.Background()

    // Create client with API token
    client := kanboard.NewClient("https://kanboard.example.com").
        WithAPIToken("your-api-token").
        WithTimeout(30 * time.Second)

    // Create a task using the fluent API
    task, err := client.Board(1).CreateTaskFromParams(ctx,
        kanboard.NewTask("My Task").
            WithDescription("Task description").
            WithPriority(2).
            WithTags("urgent", "backend"))
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created task: %d\n", task.ID)
}
```

## Authentication

The library supports two authentication methods:

### API Token (Recommended)

```go
client := kanboard.NewClient("https://kanboard.example.com").
    WithAPIToken("your-api-token")
```

### Username/Password

```go
client := kanboard.NewClient("https://kanboard.example.com").
    WithBasicAuth("admin", "password")
```

## API Usage

### Fluent API

The fluent API provides chainable, scoped operations:

```go
// Board-scoped operations
board := client.Board(projectID)
tasks, _ := board.GetTasks(ctx, kanboard.StatusActive)
columns, _ := board.GetColumns(ctx)
categories, _ := board.GetCategories(ctx)

// Task-scoped operations
task := client.Task(taskID)
task.Close(ctx)
task.MoveToNextColumn(ctx)
task.AddTag(ctx, "reviewed")
task.AddComment(ctx, userID, "Comment text")
```

### Direct API

For lower-level access, use the direct API methods:

```go
// Projects
projects, _ := client.GetAllProjects(ctx)
project, _ := client.GetProjectByID(ctx, projectID)
project, _ := client.GetProjectByName(ctx, "My Project")

// Tasks
tasks, _ := client.GetAllTasks(ctx, projectID, kanboard.StatusActive)
task, _ := client.GetTask(ctx, taskID)
task, _ := client.CreateTask(ctx, kanboard.CreateTaskRequest{
    Title:     "New Task",
    ProjectID: projectID,
})

// Global search across all projects
results, _ := client.SearchTasksGlobal(ctx, "search query")
```

### Task Creation

Use `TaskParams` for fluent task creation:

```go
params := kanboard.NewTask("Task Title").
    WithDescription("Detailed description").
    WithCategory(categoryID).
    WithOwner(userID).
    WithColor("red").
    WithPriority(3).
    WithScore(5).
    WithDueDate(time.Now().Add(7 * 24 * time.Hour)).
    WithTags("feature", "v2.0").
    InColumn(columnID).
    InSwimlane(swimlaneID)

task, err := client.Board(projectID).CreateTaskFromParams(ctx, params)
```

### Task Updates

Use `TaskUpdateParams` for partial updates:

```go
updates := kanboard.NewTaskUpdate().
    SetTitle("Updated Title").
    SetDescription("New description").
    SetPriority(1)

err := client.Task(taskID).Update(ctx, updates)
```

### Task Movement

```go
// Move to next/previous column in workflow
client.Task(taskID).MoveToNextColumn(ctx)
client.Task(taskID).MoveToPreviousColumn(ctx)

// Move to specific column
client.Task(taskID).MoveToColumn(ctx, columnID)

// Move to different project
client.Task(taskID).MoveToProject(ctx, newProjectID)
```

### Comments

```go
// Get all comments
comments, _ := client.Task(taskID).GetComments(ctx)

// Add a comment
comment, _ := client.Task(taskID).AddComment(ctx, userID, "Comment text")
```

### Links

```go
// Get task links
links, _ := client.Task(taskID).GetLinks(ctx)

// Link tasks together
client.Task(taskID).LinkTo(ctx, otherTaskID, linkID)
```

### Files

```go
// Get attached files
files, _ := client.Task(taskID).GetFiles(ctx)

// Upload a file
content := []byte("file content")
fileID, _ := client.Task(taskID).UploadFile(ctx, "document.txt", content)
```

## Error Handling

The library provides typed errors for common scenarios:

```go
task, err := client.GetTask(ctx, taskID)
if err != nil {
    if kanboard.IsNotFound(err) {
        // Handle not found
    }
    if kanboard.IsUnauthorized(err) {
        // Handle auth failure
    }
    if kanboard.IsAPIError(err) {
        // Handle Kanboard API error
    }
}
```

Available sentinel errors:
- `ErrNotFound`, `ErrProjectNotFound`, `ErrTaskNotFound`, `ErrColumnNotFound`
- `ErrUnauthorized`, `ErrForbidden`
- `ErrConnectionFailed`, `ErrTimeout`
- `ErrAlreadyInLastColumn`, `ErrAlreadyInFirstColumn`
- `ErrEmptyTitle`, `ErrInvalidProjectID`

## Thread Safety

The client is safe for concurrent use by multiple goroutines. Request IDs are generated atomically.

## Tag Operations Warning

Kanboard's `setTaskTags` API **replaces all tags**. The library implements read-modify-write internally for `AddTag` and `RemoveTag` operations. This is not atomic - concurrent tag modifications may cause data loss.

```go
// Safe: single operation
task.SetTags(ctx, "tag1", "tag2", "tag3")

// Caution: concurrent calls to AddTag/RemoveTag may conflict
task.AddTag(ctx, "new-tag")
task.RemoveTag(ctx, "old-tag")
```

## Client Configuration

```go
client := kanboard.NewClient("https://kanboard.example.com").
    WithAPIToken("token").
    WithTimeout(60 * time.Second).
    WithHTTPClient(customHTTPClient).
    WithLogger(slog.Default())
```

## License

MIT License - see [LICENSE](LICENSE) for details.
