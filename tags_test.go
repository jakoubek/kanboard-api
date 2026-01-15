package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetTaskTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getTaskTags" {
			t.Errorf("expected method=getTaskTags, got %s", req.Method)
		}

		params := req.Params.(map[string]interface{})
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		// Kanboard returns map[string]string
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"1": "urgent", "2": "backend", "5": "bug"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tags, err := client.GetTaskTags(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}
	if tags[1] != "urgent" {
		t.Errorf("expected tags[1]='urgent', got %s", tags[1])
	}
	if tags[2] != "backend" {
		t.Errorf("expected tags[2]='backend', got %s", tags[2])
	}
	if tags[5] != "bug" {
		t.Errorf("expected tags[5]='bug', got %s", tags[5])
	}
}

func TestClient_GetTaskTags_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tags, err := client.GetTaskTags(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestClient_SetTaskTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "setTaskTags" {
			t.Errorf("expected method=setTaskTags, got %s", req.Method)
		}

		params := req.Params.(map[string]interface{})
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		tags := params["tags"].([]interface{})
		if len(tags) != 2 {
			t.Errorf("expected 2 tags, got %d", len(tags))
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

	err := client.SetTaskTags(context.Background(), 1, 42, []string{"urgent", "backend"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_SetTaskTags_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		params := req.Params.(map[string]interface{})
		tags := params["tags"].([]interface{})
		if len(tags) != 0 {
			t.Errorf("expected 0 tags, got %d", len(tags))
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

	// Clear all tags by passing empty slice
	err := client.SetTaskTags(context.Background(), 1, 42, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_GetAllTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllTags" {
			t.Errorf("expected method=getAllTags, got %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "name": "urgent", "project_id": "1", "color_id": "red"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tags, err := client.GetAllTags(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Name != "urgent" {
		t.Errorf("expected name='urgent', got %s", tags[0].Name)
	}
}

func TestClient_GetTagsByProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getTagsByProject" {
			t.Errorf("expected method=getTagsByProject, got %s", req.Method)
		}

		params := req.Params.(map[string]interface{})
		if params["project_id"].(float64) != 5 {
			t.Errorf("expected project_id=5, got %v", params["project_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "name": "urgent", "project_id": "5", "color_id": "red"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tags, err := client.GetTagsByProject(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestClient_CreateTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "createTag" {
			t.Errorf("expected method=createTag, got %s", req.Method)
		}

		params := req.Params.(map[string]interface{})
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["tag"] != "new-tag" {
			t.Errorf("expected tag='new-tag', got %v", params["tag"])
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
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tagID, err := client.CreateTag(context.Background(), 1, "new-tag", "blue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tagID != 42 {
		t.Errorf("expected tagID=42, got %d", tagID)
	}
}

func TestClient_CreateTag_NoColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		params := req.Params.(map[string]interface{})
		if _, exists := params["color_id"]; exists {
			t.Error("color_id should not be present when empty")
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`42`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.CreateTag(context.Background(), 1, "new-tag", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_UpdateTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "updateTag" {
			t.Errorf("expected method=updateTag, got %s", req.Method)
		}

		params := req.Params.(map[string]interface{})
		if params["tag_id"].(float64) != 5 {
			t.Errorf("expected tag_id=5, got %v", params["tag_id"])
		}
		if params["tag"] != "updated-name" {
			t.Errorf("expected tag='updated-name', got %v", params["tag"])
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

	err := client.UpdateTag(context.Background(), 5, "updated-name", "green")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_RemoveTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "removeTag" {
			t.Errorf("expected method=removeTag, got %s", req.Method)
		}

		params := req.Params.(map[string]interface{})
		if params["tag_id"].(float64) != 5 {
			t.Errorf("expected tag_id=5, got %v", params["tag_id"])
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

	err := client.RemoveTag(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
