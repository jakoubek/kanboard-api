package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTaskUpdate(t *testing.T) {
	params := NewTaskUpdate()

	if params == nil {
		t.Fatal("expected non-nil TaskUpdateParams")
	}
	if params.title != nil {
		t.Error("expected title to be nil")
	}
}

func TestTaskUpdateParams_Chaining(t *testing.T) {
	dueDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	params := NewTaskUpdate().
		SetTitle("Updated Title").
		SetDescription("New description").
		SetColor("red").
		SetOwner(10).
		SetCategory(5).
		SetPriority(2).
		SetDueDate(dueDate).
		SetReference("JIRA-456").
		SetTags("urgent", "review")

	if *params.title != "Updated Title" {
		t.Errorf("expected title='Updated Title', got %s", *params.title)
	}
	if *params.description != "New description" {
		t.Errorf("expected description='New description', got %s", *params.description)
	}
	if *params.colorID != "red" {
		t.Errorf("expected colorID='red', got %s", *params.colorID)
	}
	if *params.ownerID != 10 {
		t.Errorf("expected ownerID=10, got %d", *params.ownerID)
	}
	if *params.categoryID != 5 {
		t.Errorf("expected categoryID=5, got %d", *params.categoryID)
	}
	if *params.priority != 2 {
		t.Errorf("expected priority=2, got %d", *params.priority)
	}
	if *params.dueDate != dueDate.Format(kanboardDateTimeFormat) {
		t.Errorf("expected dueDate=%s, got %s", dueDate.Format(kanboardDateTimeFormat), *params.dueDate)
	}
	if *params.reference != "JIRA-456" {
		t.Errorf("expected reference='JIRA-456', got %s", *params.reference)
	}
	if len(params.tags) != 2 || params.tags[0] != "urgent" {
		t.Errorf("expected tags=['urgent', 'review'], got %v", params.tags)
	}
}

func TestTaskUpdateParams_SetScore(t *testing.T) {
	params := NewTaskUpdate().SetScore(5)

	if *params.score != 5 {
		t.Errorf("expected score=5, got %d", *params.score)
	}
}

func TestTaskUpdateParams_SetStartDate(t *testing.T) {
	startDate := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	params := NewTaskUpdate().SetStartDate(startDate)

	if *params.startDate != startDate.Format(kanboardDateTimeFormat) {
		t.Errorf("expected startDate=%s, got %s", startDate.Format(kanboardDateTimeFormat), *params.startDate)
	}
}

func TestTaskUpdateParams_ClearDueDate(t *testing.T) {
	params := NewTaskUpdate().ClearDueDate()

	if params.dueDate == nil {
		t.Fatal("expected dueDate to be set")
	}
	if *params.dueDate != "" {
		t.Errorf("expected dueDate='', got %s", *params.dueDate)
	}
}

func TestTaskUpdateParams_ClearStartDate(t *testing.T) {
	params := NewTaskUpdate().ClearStartDate()

	if *params.startDate != "" {
		t.Errorf("expected startDate='', got %s", *params.startDate)
	}
}

func TestTaskUpdateParams_ClearOwner(t *testing.T) {
	params := NewTaskUpdate().ClearOwner()

	if *params.ownerID != 0 {
		t.Errorf("expected ownerID=0, got %d", *params.ownerID)
	}
}

func TestTaskUpdateParams_ClearCategory(t *testing.T) {
	params := NewTaskUpdate().ClearCategory()

	if *params.categoryID != 0 {
		t.Errorf("expected categoryID=0, got %d", *params.categoryID)
	}
}

func TestTaskUpdateParams_toUpdateTaskRequest(t *testing.T) {
	dueDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	params := NewTaskUpdate().
		SetTitle("Updated Title").
		SetDescription("Details").
		SetPriority(2).
		SetDueDate(dueDate).
		SetTags("urgent")

	req := params.toUpdateTaskRequest(42)

	if req.ID != 42 {
		t.Errorf("expected ID=42, got %d", req.ID)
	}
	if *req.Title != "Updated Title" {
		t.Errorf("expected Title='Updated Title', got %s", *req.Title)
	}
	if *req.Description != "Details" {
		t.Errorf("expected Description='Details', got %s", *req.Description)
	}
	if *req.Priority != 2 {
		t.Errorf("expected Priority=2, got %d", *req.Priority)
	}
	if *req.DateDue != dueDate.Format(kanboardDateTimeFormat) {
		t.Errorf("expected DateDue=%s, got %s", dueDate.Format(kanboardDateTimeFormat), *req.DateDue)
	}
	if len(req.Tags) != 1 || req.Tags[0] != "urgent" {
		t.Errorf("expected Tags=['urgent'], got %v", req.Tags)
	}
}

func TestTaskUpdateParams_toUpdateTaskRequest_PartialUpdate(t *testing.T) {
	// Only set title, nothing else
	params := NewTaskUpdate().SetTitle("New Title")

	req := params.toUpdateTaskRequest(42)

	if *req.Title != "New Title" {
		t.Errorf("expected Title='New Title', got %s", *req.Title)
	}
	// Unset fields should be nil
	if req.Description != nil {
		t.Error("expected Description to be nil")
	}
	if req.Priority != nil {
		t.Error("expected Priority to be nil")
	}
	if req.OwnerID != nil {
		t.Error("expected OwnerID to be nil")
	}
	if req.Tags != nil {
		t.Error("expected Tags to be nil")
	}
}

func TestTaskUpdateParams_SetTags_Empty(t *testing.T) {
	// Explicitly clearing tags
	params := NewTaskUpdate().SetTags()

	if !params.tagsSet {
		t.Error("expected tagsSet to be true")
	}
	if len(params.tags) != 0 {
		t.Errorf("expected empty tags, got %v", params.tags)
	}

	req := params.toUpdateTaskRequest(42)
	if req.Tags == nil {
		t.Error("expected Tags to be non-nil (empty slice)")
	}
	if len(req.Tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(req.Tags))
	}
}

func TestClient_UpdateTaskFromParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "updateTask" {
			t.Errorf("expected method=updateTask, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["id"].(float64) != 42 {
			t.Errorf("expected id=42, got %v", params["id"])
		}
		if params["title"] != "Updated via Params" {
			t.Errorf("expected title='Updated via Params', got %v", params["title"])
		}
		if params["priority"].(float64) != 3 {
			t.Errorf("expected priority=3, got %v", params["priority"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	params := NewTaskUpdate().
		SetTitle("Updated via Params").
		SetPriority(3)

	err := client.UpdateTaskFromParams(context.Background(), 42, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_UpdateTaskFromParams_OnlySetFieldsSent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		params := req.Params.(map[string]any)

		// Only title should be present, not priority or other fields
		if _, hasTitle := params["title"]; !hasTitle {
			t.Error("expected title to be present")
		}
		if _, hasPriority := params["priority"]; hasPriority {
			t.Error("expected priority to NOT be present")
		}
		if _, hasOwner := params["owner_id"]; hasOwner {
			t.Error("expected owner_id to NOT be present")
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	// Only set title
	params := NewTaskUpdate().SetTitle("Only Title")

	err := client.UpdateTaskFromParams(context.Background(), 42, params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
