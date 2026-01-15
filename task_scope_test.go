package kanboard

import (
	"context"
	"encoding/json"
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
