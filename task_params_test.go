package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	params := NewTask("Test Task")

	if params.title != "Test Task" {
		t.Errorf("expected title='Test Task', got %s", params.title)
	}
	if params.description != nil {
		t.Error("expected description to be nil")
	}
	if params.columnID != nil {
		t.Error("expected columnID to be nil")
	}
}

func TestTaskParams_Chaining(t *testing.T) {
	dueDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	params := NewTask("My Task").
		WithDescription("Task details").
		InColumn(2).
		WithCategory(5).
		WithOwner(10).
		WithColor("blue").
		WithPriority(2).
		WithDueDate(dueDate).
		WithTags("urgent", "backend").
		WithReference("JIRA-123")

	if params.title != "My Task" {
		t.Errorf("expected title='My Task', got %s", params.title)
	}
	if *params.description != "Task details" {
		t.Errorf("expected description='Task details', got %s", *params.description)
	}
	if *params.columnID != 2 {
		t.Errorf("expected columnID=2, got %d", *params.columnID)
	}
	if *params.categoryID != 5 {
		t.Errorf("expected categoryID=5, got %d", *params.categoryID)
	}
	if *params.ownerID != 10 {
		t.Errorf("expected ownerID=10, got %d", *params.ownerID)
	}
	if *params.colorID != "blue" {
		t.Errorf("expected colorID='blue', got %s", *params.colorID)
	}
	if *params.priority != 2 {
		t.Errorf("expected priority=2, got %d", *params.priority)
	}
	if *params.dueDate != dueDate.Format(kanboardDateTimeFormat) {
		t.Errorf("expected dueDate=%s, got %s", dueDate.Format(kanboardDateTimeFormat), *params.dueDate)
	}
	if len(params.tags) != 2 || params.tags[0] != "urgent" || params.tags[1] != "backend" {
		t.Errorf("expected tags=['urgent', 'backend'], got %v", params.tags)
	}
	if *params.reference != "JIRA-123" {
		t.Errorf("expected reference='JIRA-123', got %s", *params.reference)
	}
}

func TestTaskParams_WithScore(t *testing.T) {
	params := NewTask("Task").WithScore(8)

	if *params.score != 8 {
		t.Errorf("expected score=8, got %d", *params.score)
	}
}

func TestTaskParams_WithCreator(t *testing.T) {
	params := NewTask("Task").WithCreator(5)

	if *params.creatorID != 5 {
		t.Errorf("expected creatorID=5, got %d", *params.creatorID)
	}
}

func TestTaskParams_WithStartDate(t *testing.T) {
	startDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	params := NewTask("Task").WithStartDate(startDate)

	if *params.startDate != startDate.Format(kanboardDateTimeFormat) {
		t.Errorf("expected startDate=%s, got %s", startDate.Format(kanboardDateTimeFormat), *params.startDate)
	}
}

func TestTaskParams_InSwimlane(t *testing.T) {
	params := NewTask("Task").InSwimlane(3)

	if *params.swimlaneID != 3 {
		t.Errorf("expected swimlaneID=3, got %d", *params.swimlaneID)
	}
}

func TestTaskParams_toCreateTaskRequest(t *testing.T) {
	dueDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	params := NewTask("Test Task").
		WithDescription("Details").
		InColumn(2).
		WithCategory(5).
		WithOwner(10).
		WithColor("blue").
		WithPriority(2).
		WithScore(8).
		WithDueDate(dueDate).
		WithTags("urgent", "backend").
		WithReference("REF-123")

	req := params.toCreateTaskRequest(42)

	if req.Title != "Test Task" {
		t.Errorf("expected Title='Test Task', got %s", req.Title)
	}
	if req.ProjectID != 42 {
		t.Errorf("expected ProjectID=42, got %d", req.ProjectID)
	}
	if req.Description != "Details" {
		t.Errorf("expected Description='Details', got %s", req.Description)
	}
	if req.ColumnID != 2 {
		t.Errorf("expected ColumnID=2, got %d", req.ColumnID)
	}
	if req.CategoryID != 5 {
		t.Errorf("expected CategoryID=5, got %d", req.CategoryID)
	}
	if req.OwnerID != 10 {
		t.Errorf("expected OwnerID=10, got %d", req.OwnerID)
	}
	if req.ColorID != "blue" {
		t.Errorf("expected ColorID='blue', got %s", req.ColorID)
	}
	if req.Priority != 2 {
		t.Errorf("expected Priority=2, got %d", req.Priority)
	}
	if req.Score != 8 {
		t.Errorf("expected Score=8, got %d", req.Score)
	}
	if req.DateDue != dueDate.Format(kanboardDateTimeFormat) {
		t.Errorf("expected DateDue=%s, got %s", dueDate.Format(kanboardDateTimeFormat), req.DateDue)
	}
	if len(req.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(req.Tags))
	}
	if req.Reference != "REF-123" {
		t.Errorf("expected Reference='REF-123', got %s", req.Reference)
	}
}

func TestTaskParams_toCreateTaskRequest_MinimalParams(t *testing.T) {
	params := NewTask("Simple Task")

	req := params.toCreateTaskRequest(1)

	if req.Title != "Simple Task" {
		t.Errorf("expected Title='Simple Task', got %s", req.Title)
	}
	if req.ProjectID != 1 {
		t.Errorf("expected ProjectID=1, got %d", req.ProjectID)
	}
	// Unset fields should have zero values
	if req.Description != "" {
		t.Errorf("expected Description='', got %s", req.Description)
	}
	if req.ColumnID != 0 {
		t.Errorf("expected ColumnID=0, got %d", req.ColumnID)
	}
	if len(req.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(req.Tags))
	}
}

func TestBoardScope_CreateTaskFromParams(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			// First call: createTask
			if req.Method != "createTask" {
				t.Errorf("expected method=createTask, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			if params["title"] != "Fluent Task" {
				t.Errorf("expected title='Fluent Task', got %v", params["title"])
			}
			if params["project_id"].(float64) != 1 {
				t.Errorf("expected project_id=1, got %v", params["project_id"])
			}
			if params["description"] != "Task details" {
				t.Errorf("expected description='Task details', got %v", params["description"])
			}
			if params["color_id"] != "red" {
				t.Errorf("expected color_id='red', got %v", params["color_id"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`42`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Fluent Task", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	params := NewTask("Fluent Task").
		WithDescription("Task details").
		WithColor("red")

	task, err := client.Board(1).CreateTaskFromParams(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(task.ID) != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
}
