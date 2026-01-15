package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetAllTaskLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllTaskLinks" {
			t.Errorf("expected method=getAllTaskLinks, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "1", "link_id": "2", "task_id": "42", "opposite_task_id": "100", "label": "blocks", "title": "Blocked Task"},
				{"id": "2", "link_id": "3", "task_id": "42", "opposite_task_id": "101", "label": "is duplicated by", "title": "Duplicate Task"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	links, err := client.GetAllTaskLinks(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(links) != 2 {
		t.Errorf("expected 2 links, got %d", len(links))
	}
	if links[0].Label != "blocks" {
		t.Errorf("expected label='blocks', got %s", links[0].Label)
	}
	if int(links[0].OppositeTaskID) != 100 {
		t.Errorf("expected opposite_task_id=100, got %d", links[0].OppositeTaskID)
	}
}

func TestClient_GetAllTaskLinks_Empty(t *testing.T) {
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

	links, err := client.GetAllTaskLinks(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(links) != 0 {
		t.Errorf("expected 0 links, got %d", len(links))
	}
}

func TestClient_CreateTaskLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "createTaskLink" {
			t.Errorf("expected method=createTaskLink, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}
		if params["opposite_task_id"].(float64) != 100 {
			t.Errorf("expected opposite_task_id=100, got %v", params["opposite_task_id"])
		}
		if params["link_id"].(float64) != 2 {
			t.Errorf("expected link_id=2, got %v", params["link_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`10`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	linkID, err := client.CreateTaskLink(context.Background(), 42, 100, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if linkID != 10 {
		t.Errorf("expected linkID=10, got %d", linkID)
	}
}

func TestClient_CreateTaskLink_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`0`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.CreateTaskLink(context.Background(), 42, 100, 2)
	if err == nil {
		t.Fatal("expected error for failed link creation")
	}
}

func TestClient_RemoveTaskLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "removeTaskLink" {
			t.Errorf("expected method=removeTaskLink, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_link_id"].(float64) != 5 {
			t.Errorf("expected task_link_id=5, got %v", params["task_link_id"])
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

	err := client.RemoveTaskLink(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_RemoveTaskLink_Failure(t *testing.T) {
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

	err := client.RemoveTaskLink(context.Background(), 5)
	if err == nil {
		t.Fatal("expected error for failed delete")
	}
}

func TestTaskScope_GetLinks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllTaskLinks" {
			t.Errorf("expected method=getAllTaskLinks, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "link_id": "2", "task_id": "42", "opposite_task_id": "100", "label": "blocks"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	links, err := client.Task(42).GetLinks(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(links) != 1 {
		t.Errorf("expected 1 link, got %d", len(links))
	}
}

func TestTaskScope_LinkTo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "createTaskLink" {
			t.Errorf("expected method=createTaskLink, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}
		if params["opposite_task_id"].(float64) != 100 {
			t.Errorf("expected opposite_task_id=100, got %v", params["opposite_task_id"])
		}
		if params["link_id"].(float64) != 2 {
			t.Errorf("expected link_id=2, got %v", params["link_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`10`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	err := client.Task(42).LinkTo(context.Background(), 100, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
