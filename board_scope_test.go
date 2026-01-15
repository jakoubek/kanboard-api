package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Board(t *testing.T) {
	client := NewClient("http://example.com").WithAPIToken("test")

	board := client.Board(42)
	if board == nil {
		t.Fatal("expected non-nil BoardScope")
	}
	if board.projectID != 42 {
		t.Errorf("expected projectID=42, got %d", board.projectID)
	}
	if board.client != client {
		t.Error("expected BoardScope to reference the same client")
	}
}

func TestBoardScope_GetColumns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getColumns" {
			t.Errorf("expected method=getColumns, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "title": "Backlog", "position": "1", "project_id": "1"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	columns, err := client.Board(1).GetColumns(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 1 {
		t.Errorf("expected 1 column, got %d", len(columns))
	}
	if columns[0].Title != "Backlog" {
		t.Errorf("expected title='Backlog', got %s", columns[0].Title)
	}
}

func TestBoardScope_GetCategories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllCategories" {
			t.Errorf("expected method=getAllCategories, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "name": "Bug", "project_id": "1"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	categories, err := client.Board(1).GetCategories(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(categories) != 1 {
		t.Errorf("expected 1 category, got %d", len(categories))
	}
	if categories[0].Name != "Bug" {
		t.Errorf("expected name='Bug', got %s", categories[0].Name)
	}
}

func TestBoardScope_GetTasks(t *testing.T) {
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
			t.Errorf("expected status_id=1 (active), got %v", params["status_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "title": "Task One", "project_id": "1", "is_active": "1"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.Board(1).GetTasks(context.Background(), StatusActive)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Title != "Task One" {
		t.Errorf("expected title='Task One', got %s", tasks[0].Title)
	}
}

func TestBoardScope_SearchTasks(t *testing.T) {
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
		if params["query"] != "status:open" {
			t.Errorf("expected query='status:open', got %v", params["query"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "title": "Task One", "project_id": "1", "is_active": "1"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tasks, err := client.Board(1).SearchTasks(context.Background(), "status:open")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasks))
	}
}

func TestBoardScope_CreateTask(t *testing.T) {
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
			// BoardScope should set the project_id
			if params["project_id"].(float64) != 1 {
				t.Errorf("expected project_id=1, got %v", params["project_id"])
			}
			if params["title"] != "New Task" {
				t.Errorf("expected title='New Task', got %v", params["title"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`42`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: getTask to fetch created task
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

	task, err := client.Board(1).CreateTask(context.Background(), CreateTaskRequest{
		Title: "New Task",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(task.ID) != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if int(task.ProjectID) != 1 {
		t.Errorf("expected ProjectID=1, got %d", task.ProjectID)
	}
}

func TestBoardScope_CreateTask_OverridesProjectID(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			params := req.Params.(map[string]any)
			// Even if user provides different project_id, BoardScope should override
			if params["project_id"].(float64) != 5 {
				t.Errorf("expected project_id=5 (from BoardScope), got %v", params["project_id"])
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
				Result:  json.RawMessage(`{"id": "42", "title": "New Task", "project_id": "5", "column_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	// User tries to set different project_id, but BoardScope should override
	_, err := client.Board(5).CreateTask(context.Background(), CreateTaskRequest{
		Title:     "New Task",
		ProjectID: 999, // This should be overridden to 5
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
