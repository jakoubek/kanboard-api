package kanboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetAllProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllProjects" {
			t.Errorf("expected method=getAllProjects, got %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "1", "name": "Project One", "is_active": "1"},
				{"id": "2", "name": "Project Two", "is_active": "0"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	projects, err := client.GetAllProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(projects))
	}
	if projects[0].Name != "Project One" {
		t.Errorf("expected name='Project One', got %s", projects[0].Name)
	}
	if int(projects[0].ID) != 1 {
		t.Errorf("expected ID=1, got %d", projects[0].ID)
	}
}

func TestClient_GetAllProjects_Empty(t *testing.T) {
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

	projects, err := client.GetAllProjects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("expected 0 projects, got %d", len(projects))
	}
}

func TestClient_GetProjectByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getProjectById" {
			t.Errorf("expected method=getProjectById, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 42 {
			t.Errorf("expected project_id=42, got %v", params["project_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "42", "name": "Test Project", "is_active": "1"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	project, err := client.GetProjectByID(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(project.ID) != 42 {
		t.Errorf("expected ID=42, got %d", project.ID)
	}
	if project.Name != "Test Project" {
		t.Errorf("expected name='Test Project', got %s", project.Name)
	}
}

func TestClient_GetProjectByID_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Kanboard returns null for non-existent projects
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`null`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.GetProjectByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent project")
	}

	if !errors.Is(err, ErrProjectNotFound) {
		t.Errorf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestClient_GetProjectByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getProjectByName" {
			t.Errorf("expected method=getProjectByName, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["name"] != "Test Project" {
			t.Errorf("expected name='Test Project', got %v", params["name"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "42", "name": "Test Project", "is_active": "1"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	project, err := client.GetProjectByName(context.Background(), "Test Project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if project.Name != "Test Project" {
		t.Errorf("expected name='Test Project', got %s", project.Name)
	}
}

func TestClient_GetProjectByName_NotFound(t *testing.T) {
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

	_, err := client.GetProjectByName(context.Background(), "Non-Existent")
	if err == nil {
		t.Fatal("expected error for non-existent project")
	}

	if !errors.Is(err, ErrProjectNotFound) {
		t.Errorf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestClient_GetProjectByID_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		select {}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.GetProjectByID(ctx, 42)
	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
}
