package kanboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Task(t *testing.T) {
	client := NewClient("http://example.com").WithAPIToken("test")

	scope := client.Task(42)
	if scope == nil {
		t.Fatal("expected non-nil TaskScope")
	}
	if scope.taskID != 42 {
		t.Errorf("expected taskID=42, got %d", scope.taskID)
	}
	if scope.client != client {
		t.Error("expected TaskScope to reference the same client")
	}
}

func TestTaskScope_Get(t *testing.T) {
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
			Result:  json.RawMessage(`{"id": "42", "title": "Test Task", "project_id": "1", "is_active": "1"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	task, err := client.Task(42).Get(context.Background())
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

func TestTaskScope_Close(t *testing.T) {
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

	err := client.Task(42).Close(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_Open(t *testing.T) {
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

	err := client.Task(42).Open(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_MoveToColumn(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			// First call: getTask to get project_id and swimlane_id
			if req.Method != "getTask" {
				t.Errorf("expected method=getTask, got %s", req.Method)
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test Task", "project_id": "1", "column_id": "1", "swimlane_id": "0", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: moveTaskPosition
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
			if params["column_id"].(float64) != 3 {
				t.Errorf("expected column_id=3, got %v", params["column_id"])
			}
			if params["position"].(float64) != 0 {
				t.Errorf("expected position=0, got %v", params["position"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).MoveToColumn(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_MoveToNextColumn(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// First call: getTask
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "project_id": "1", "column_id": "2", "swimlane_id": "0", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// Second call: getColumns
			if req.Method != "getColumns" {
				t.Errorf("expected method=getColumns, got %s", req.Method)
			}
			// Return columns in position order
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`[
					{"id": "1", "title": "Backlog", "position": "1", "project_id": "1"},
					{"id": "2", "title": "In Progress", "position": "2", "project_id": "1"},
					{"id": "3", "title": "Done", "position": "3", "project_id": "1"}
				]`),
			}
			json.NewEncoder(w).Encode(resp)
		case 3:
			// Third call: moveTaskPosition
			if req.Method != "moveTaskPosition" {
				t.Errorf("expected method=moveTaskPosition, got %s", req.Method)
			}
			params := req.Params.(map[string]any)
			if params["column_id"].(float64) != 3 {
				t.Errorf("expected column_id=3 (Done), got %v", params["column_id"])
			}
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).MoveToNextColumn(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_MoveToNextColumn_AlreadyInLastColumn(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// Task is already in column 3 (Done)
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "project_id": "1", "column_id": "3", "swimlane_id": "0", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// getColumns
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`[
					{"id": "1", "title": "Backlog", "position": "1", "project_id": "1"},
					{"id": "2", "title": "In Progress", "position": "2", "project_id": "1"},
					{"id": "3", "title": "Done", "position": "3", "project_id": "1"}
				]`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Error("moveTaskPosition should not be called when already in last column")
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).MoveToNextColumn(context.Background())
	if err == nil {
		t.Fatal("expected error for task already in last column")
	}

	if !errors.Is(err, ErrAlreadyInLastColumn) {
		t.Errorf("expected ErrAlreadyInLastColumn, got %v", err)
	}
}

func TestTaskScope_MoveToPreviousColumn(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// Task in column 2 (In Progress)
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "project_id": "1", "column_id": "2", "swimlane_id": "0", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// getColumns
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`[
					{"id": "1", "title": "Backlog", "position": "1", "project_id": "1"},
					{"id": "2", "title": "In Progress", "position": "2", "project_id": "1"},
					{"id": "3", "title": "Done", "position": "3", "project_id": "1"}
				]`),
			}
			json.NewEncoder(w).Encode(resp)
		case 3:
			// moveTaskPosition
			params := req.Params.(map[string]any)
			if params["column_id"].(float64) != 1 {
				t.Errorf("expected column_id=1 (Backlog), got %v", params["column_id"])
			}
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).MoveToPreviousColumn(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_MoveToPreviousColumn_AlreadyInFirstColumn(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// Task is already in column 1 (Backlog)
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "project_id": "1", "column_id": "1", "swimlane_id": "0", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// getColumns
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: json.RawMessage(`[
					{"id": "1", "title": "Backlog", "position": "1", "project_id": "1"},
					{"id": "2", "title": "In Progress", "position": "2", "project_id": "1"},
					{"id": "3", "title": "Done", "position": "3", "project_id": "1"}
				]`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Error("moveTaskPosition should not be called when already in first column")
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).MoveToPreviousColumn(context.Background())
	if err == nil {
		t.Fatal("expected error for task already in first column")
	}

	if !errors.Is(err, ErrAlreadyInFirstColumn) {
		t.Errorf("expected ErrAlreadyInFirstColumn, got %v", err)
	}
}

func TestTaskScope_MoveToProject(t *testing.T) {
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

	err := client.Task(42).MoveToProject(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_Update(t *testing.T) {
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

	err := client.Task(42).Update(context.Background(), NewTaskUpdate().SetTitle("Updated Title"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_Chaining(t *testing.T) {
	// Verify we can chain operations on the same task
	client := NewClient("http://example.com").WithAPIToken("test")

	// These should all return the same underlying taskID
	scope1 := client.Task(42)
	scope2 := client.Task(42)

	if scope1.taskID != scope2.taskID {
		t.Error("expected same taskID for same task scope")
	}
}

func TestTaskScope_GetTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getTaskTags" {
			t.Errorf("expected method=getTaskTags, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"1": "urgent", "2": "backend"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tags, err := client.Task(42).GetTags(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
	if tags[1] != "urgent" {
		t.Errorf("expected tags[1]='urgent', got %s", tags[1])
	}
}

func TestTaskScope_SetTags(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			// First call: getTask to get project_id
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: setTaskTags
			if req.Method != "setTaskTags" {
				t.Errorf("expected method=setTaskTags, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			if params["project_id"].(float64) != 1 {
				t.Errorf("expected project_id=1, got %v", params["project_id"])
			}
			if params["task_id"].(float64) != 42 {
				t.Errorf("expected task_id=42, got %v", params["task_id"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).SetTags(context.Background(), "urgent", "review")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_ClearTags(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			// First call: getTask
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: setTaskTags with empty array
			if req.Method != "setTaskTags" {
				t.Errorf("expected method=setTaskTags, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			tags, ok := params["tags"].([]any)
			if !ok {
				// Tags might be nil if passed as nil slice
				tags = []any{}
			}
			if len(tags) != 0 {
				t.Errorf("expected empty tags array, got %v", tags)
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).ClearTags(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_AddTag(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// First call: getTask
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// Second call: getTaskTags
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"1": "existing"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 3:
			// Third call: setTaskTags
			if req.Method != "setTaskTags" {
				t.Errorf("expected method=setTaskTags, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			tags := params["tags"].([]any)
			if len(tags) != 2 {
				t.Errorf("expected 2 tags, got %d", len(tags))
			}
			// Check new tag is present
			hasNew := false
			for _, tag := range tags {
				if tag == "new-tag" {
					hasNew = true
				}
			}
			if !hasNew {
				t.Error("expected 'new-tag' in tags")
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).AddTag(context.Background(), "new-tag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_AddTag_Idempotent(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// First call: getTask
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// Second call: getTaskTags - tag already exists
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"1": "existing-tag"}`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			// Should not reach setTaskTags
			t.Error("setTaskTags should not be called when tag already exists")
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	// Adding tag that already exists - should be no-op
	err := client.Task(42).AddTag(context.Background(), "existing-tag")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_RemoveTag(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// First call: getTask
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// Second call: getTaskTags
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"1": "keep", "2": "remove-me"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 3:
			// Third call: setTaskTags
			if req.Method != "setTaskTags" {
				t.Errorf("expected method=setTaskTags, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			tags := params["tags"].([]any)
			if len(tags) != 1 {
				t.Errorf("expected 1 tag after removal, got %d", len(tags))
			}
			// Check removed tag is not present
			for _, tag := range tags {
				if tag == "remove-me" {
					t.Error("'remove-me' should have been filtered out")
				}
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`true`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).RemoveTag(context.Background(), "remove-me")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_RemoveTag_Idempotent(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		switch callCount {
		case 1:
			// First call: getTask
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "42", "title": "Test", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		case 2:
			// Second call: getTaskTags - tag doesn't exist
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"1": "other-tag"}`),
			}
			json.NewEncoder(w).Encode(resp)
		default:
			// Should not reach setTaskTags
			t.Error("setTaskTags should not be called when tag doesn't exist")
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	// Removing tag that doesn't exist - should be no-op
	err := client.Task(42).RemoveTag(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTaskScope_HasTag_True(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"1": "urgent", "2": "backend"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	hasTag, err := client.Task(42).HasTag(context.Background(), "urgent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !hasTag {
		t.Error("expected HasTag to return true")
	}
}

func TestTaskScope_HasTag_False(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"1": "urgent"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	hasTag, err := client.Task(42).HasTag(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hasTag {
		t.Error("expected HasTag to return false")
	}
}
