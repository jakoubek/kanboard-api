package kanboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetColumns(t *testing.T) {
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

		// Return columns in wrong order to test sorting
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "3", "title": "Done", "position": "3", "project_id": "1", "task_limit": "0"},
				{"id": "1", "title": "Backlog", "position": "1", "project_id": "1", "task_limit": "10"},
				{"id": "2", "title": "In Progress", "position": "2", "project_id": "1", "task_limit": "5"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	columns, err := client.GetColumns(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 3 {
		t.Errorf("expected 3 columns, got %d", len(columns))
	}

	// Verify sorted by position
	if columns[0].Title != "Backlog" {
		t.Errorf("expected first column='Backlog', got %s", columns[0].Title)
	}
	if columns[1].Title != "In Progress" {
		t.Errorf("expected second column='In Progress', got %s", columns[1].Title)
	}
	if columns[2].Title != "Done" {
		t.Errorf("expected third column='Done', got %s", columns[2].Title)
	}
}

func TestClient_GetColumns_Empty(t *testing.T) {
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

	columns, err := client.GetColumns(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(columns) != 0 {
		t.Errorf("expected 0 columns, got %d", len(columns))
	}
}

func TestClient_GetColumn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getColumn" {
			t.Errorf("expected method=getColumn, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["column_id"].(float64) != 5 {
			t.Errorf("expected column_id=5, got %v", params["column_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "5", "title": "Backlog", "position": "1", "project_id": "1", "task_limit": "10"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	column, err := client.GetColumn(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(column.ID) != 5 {
		t.Errorf("expected ID=5, got %d", column.ID)
	}
	if column.Title != "Backlog" {
		t.Errorf("expected title='Backlog', got %s", column.Title)
	}
	if int(column.TaskLimit) != 10 {
		t.Errorf("expected task_limit=10, got %d", column.TaskLimit)
	}
}

func TestClient_GetColumn_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Kanboard returns null for non-existent columns
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`null`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.GetColumn(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent column")
	}

	if !errors.Is(err, ErrColumnNotFound) {
		t.Errorf("expected ErrColumnNotFound, got %v", err)
	}
}

func TestClient_GetColumns_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetColumns(ctx, 1)
	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
}
