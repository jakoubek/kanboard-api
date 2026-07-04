package kanboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_GetTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getTask" {
			t.Errorf("expected method=getTask, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "42", "title": "Test Task", "project_id": "1", "column_id": "1", "is_active": "1"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	task, err := client.GetTask(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(task.ID) != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title='Test Task', got %s", task.Title)
	}
	if !bool(task.IsActive) {
		t.Error("expected task to be active")
	}
	// Missing time fields must decode to 0 (backward compatibility).
	if float64(task.TimeEstimated) != 0 {
		t.Errorf("expected TimeEstimated=0 for missing field, got %v", task.TimeEstimated)
	}
	if float64(task.TimeSpent) != 0 {
		t.Errorf("expected TimeSpent=0 for missing field, got %v", task.TimeSpent)
	}
}

func TestClient_GetTask_TimeFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "42", "title": "Test Task", "project_id": "1", "column_id": "1", "is_active": "1", "time_estimated": "8", "time_spent": "2.5"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	task, err := client.GetTask(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if float64(task.TimeEstimated) != 8 {
		t.Errorf("expected TimeEstimated=8, got %v", task.TimeEstimated)
	}
	if float64(task.TimeSpent) != 2.5 {
		t.Errorf("expected TimeSpent=2.5, got %v", task.TimeSpent)
	}
}

func TestClient_GetTask_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`null`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.GetTask(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestClient_GetTaskByReference(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getTaskByReference" {
			t.Errorf("expected method=getTaskByReference, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["reference"] != "EXT-123" {
			t.Errorf("expected reference='EXT-123', got %v", params["reference"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "42", "title": "Test Task", "project_id": "1", "column_id": "1", "is_active": "1", "reference": "EXT-123"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	task, err := client.GetTaskByReference(context.Background(), 1, "EXT-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(task.ID) != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected title='Test Task', got %s", task.Title)
	}
}

func TestClient_GetTaskByReference_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`null`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.GetTaskByReference(context.Background(), 1, "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}

	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestClient_GetAllTasks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllTasks" {
			t.Errorf("expected method=getAllTasks, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["status_id"].(float64) != 1 {
			t.Errorf("expected status_id=1, got %v", params["status_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "1", "title": "Task One", "project_id": "1", "is_active": "1"},
				{"id": "2", "title": "Task Two", "project_id": "1", "is_active": "1"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.GetAllTasks(context.Background(), 1, StatusActive)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].Title != "Task One" {
		t.Errorf("expected first task='Task One', got %s", tasks[0].Title)
	}
}

func TestClient_GetAllTasks_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.GetAllTasks(context.Background(), 1, StatusInactive)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestClient_CreateTask(t *testing.T) {
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
			if params["title"] != "New Task" {
				t.Errorf("expected title='New Task', got %v", params["title"])
			}
			if params["project_id"].(float64) != 1 {
				t.Errorf("expected project_id=1, got %v", params["project_id"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`42`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: getTask to fetch created task
			if req.Method != "getTask" {
				t.Errorf("expected method=getTask, got %s", req.Method)
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "New Task", "project_id": "1", "column_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	task, err := client.CreateTask(context.Background(), CreateTaskRequest{
		Title:     "New Task",
		ProjectID: 1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(task.ID) != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if task.Title != "New Task" {
		t.Errorf("expected title='New Task', got %s", task.Title)
	}
}

func TestClient_CreateTask_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method == "createTask" {
			params := req.Params.(map[string]any)
			if params["description"] != "Task description" {
				t.Errorf("expected description='Task description', got %v", params["description"])
			}
			if params["color_id"] != "blue" {
				t.Errorf("expected color_id='blue', got %v", params["color_id"])
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
				Result:  json.RawMessage(`{"id": "42", "title": "New Task", "project_id": "1", "column_id": "1", "is_active": "1", "description": "Task description", "color_id": "blue"}`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.CreateTask(context.Background(), CreateTaskRequest{
		Title:       "New Task",
		ProjectID:   1,
		Description: "Task description",
		ColorID:     "blue",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_CreateTask_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Return 0 (false) to indicate failure
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`0`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.CreateTask(context.Background(), CreateTaskRequest{
		Title:     "New Task",
		ProjectID: 1,
	})
	if err == nil {
		t.Fatal("expected error for failed task creation")
	}
}

func TestClient_UpdateTask(t *testing.T) {
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
		if params["title"] != "Updated Title" {
			t.Errorf("expected title='Updated Title', got %v", params["title"])
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

	title := "Updated Title"
	err := client.UpdateTask(context.Background(), UpdateTaskRequest{
		ID:    42,
		Title: &title,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_UpdateTask_PartialUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		params := req.Params.(map[string]any)
		// Only priority should be set
		if _, hasTitle := params["title"]; hasTitle {
			t.Error("title should not be present in partial update")
		}
		if params["priority"].(float64) != 2 {
			t.Errorf("expected priority=2, got %v", params["priority"])
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

	priority := 2
	err := client.UpdateTask(context.Background(), UpdateTaskRequest{
		ID:       42,
		Priority: &priority,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_UpdateTask_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`false`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	title := "Updated Title"
	err := client.UpdateTask(context.Background(), UpdateTaskRequest{
		ID:    42,
		Title: &title,
	})
	if err == nil {
		t.Fatal("expected error for failed update")
	}
}

func TestClient_CloseTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "closeTask" {
			t.Errorf("expected method=closeTask, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
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

	err := client.CloseTask(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_CloseTask_AlreadyClosed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Kanboard returns false if task is already closed
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`false`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.CloseTask(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error for already closed task")
	}

	if !errors.Is(err, ErrTaskClosed) {
		t.Errorf("expected ErrTaskClosed, got %v", err)
	}
}

func TestClient_OpenTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "openTask" {
			t.Errorf("expected method=openTask, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
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

	err := client.OpenTask(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_OpenTask_AlreadyOpen(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Kanboard returns false if task is already open
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`false`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.OpenTask(context.Background(), 42)
	if err == nil {
		t.Fatal("expected error for already open task")
	}

	if !errors.Is(err, ErrTaskOpen) {
		t.Errorf("expected ErrTaskOpen, got %v", err)
	}
}

func TestClient_GetTask_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetTask(ctx, 42)
	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
}

func TestClient_SearchTasks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "searchTasks" {
			t.Errorf("expected method=searchTasks, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["query"] != "status:open assignee:me" {
			t.Errorf("expected query='status:open assignee:me', got %v", params["query"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "1", "title": "Task One", "project_id": "1", "is_active": "1"},
				{"id": "2", "title": "Task Two", "project_id": "1", "is_active": "1"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.SearchTasks(context.Background(), 1, "status:open assignee:me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestClient_SearchTasks_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.SearchTasks(context.Background(), 1, "title:nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestClient_MoveTaskPosition(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "moveTaskPosition" {
			t.Errorf("expected method=moveTaskPosition, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}
		if params["column_id"].(float64) != 2 {
			t.Errorf("expected column_id=2, got %v", params["column_id"])
		}
		if params["position"].(float64) != 1 {
			t.Errorf("expected position=1, got %v", params["position"])
		}
		if params["swimlane_id"].(float64) != 0 {
			t.Errorf("expected swimlane_id=0, got %v", params["swimlane_id"])
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

	err := client.MoveTaskPosition(context.Background(), 1, 42, 2, 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_MoveTaskPosition_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`false`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.MoveTaskPosition(context.Background(), 1, 42, 999, 1, 0)
	if err == nil {
		t.Fatal("expected error for failed move")
	}

	// Verify it's an OperationFailedError with helpful hints
	if !IsOperationFailed(err) {
		t.Errorf("expected OperationFailedError, got %T", err)
	}

	// Error message should contain actionable hints
	errMsg := err.Error()
	if !strings.Contains(errMsg, "moveTaskPosition") {
		t.Error("error should mention operation name")
	}
	if !strings.Contains(errMsg, "possible causes") {
		t.Error("error should include possible causes")
	}
}

func TestClient_MoveTaskToProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "moveTaskToProject" {
			t.Errorf("expected method=moveTaskToProject, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}
		if params["project_id"].(float64) != 5 {
			t.Errorf("expected project_id=5, got %v", params["project_id"])
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

	err := client.MoveTaskToProject(context.Background(), 42, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_MoveTaskToProject_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`false`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.MoveTaskToProject(context.Background(), 42, 999)
	if err == nil {
		t.Fatal("expected error for failed move")
	}
}

func TestClient_SearchTasksGlobally(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		switch req.Method {
		case "getAllProjects":
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`[
					{"id": "1", "name": "Project One", "is_active": "1"},
					{"id": "2", "name": "Project Two", "is_active": "1"}
				]`),
			}
			json.NewEncoder(w).Encode(resp)
		case "searchTasks":
			params := req.Params.(map[string]any)
			projectID := int(params["project_id"].(float64))
			query := params["query"].(string)

			if query != "status:open" {
				t.Errorf("expected query='status:open', got %s", query)
			}

			// Return different tasks for each project
			var result string
			if projectID == 1 {
				result = `[{"id": "1", "title": "Task from P1", "project_id": "1", "is_active": "1"}]`
			} else {
				result = `[{"id": "2", "title": "Task from P2", "project_id": "2", "is_active": "1"}]`
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(result),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.SearchTasksGlobally(context.Background(), "status:open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks from 2 projects, got %d", len(tasks))
	}
}

func TestClient_SearchTasksGlobally_NoProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.SearchTasksGlobally(context.Background(), "status:open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestClient_SearchTasksGlobally_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method == "getAllProjects" {
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`[{"id": "1", "name": "Project", "is_active": "1"}]`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Hang forever for searchTasks
			select {}
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.SearchTasksGlobally(ctx, "status:open")
	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
}
