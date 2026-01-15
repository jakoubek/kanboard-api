package kanboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetAllCategories(t *testing.T) {
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
			Result: json.RawMessage(`[
				{"id": "1", "name": "Bug", "project_id": "1", "color_id": "red"},
				{"id": "2", "name": "Feature", "project_id": "1", "color_id": "blue"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	categories, err := client.GetAllCategories(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}

	if categories[0].Name != "Bug" {
		t.Errorf("expected first category='Bug', got %s", categories[0].Name)
	}
	if categories[1].Name != "Feature" {
		t.Errorf("expected second category='Feature', got %s", categories[1].Name)
	}
}

func TestClient_GetAllCategories_Empty(t *testing.T) {
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

	categories, err := client.GetAllCategories(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(categories) != 0 {
		t.Errorf("expected 0 categories, got %d", len(categories))
	}
}

func TestClient_GetCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getCategory" {
			t.Errorf("expected method=getCategory, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["category_id"].(float64) != 5 {
			t.Errorf("expected category_id=5, got %v", params["category_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "5", "name": "Bug", "project_id": "1", "color_id": "red"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	category, err := client.GetCategory(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(category.ID) != 5 {
		t.Errorf("expected ID=5, got %d", category.ID)
	}
	if category.Name != "Bug" {
		t.Errorf("expected name='Bug', got %s", category.Name)
	}
	if category.ColorID != "red" {
		t.Errorf("expected color_id='red', got %s", category.ColorID)
	}
}

func TestClient_GetCategory_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Kanboard returns null for non-existent categories
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`null`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	_, err := client.GetCategory(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent category")
	}

	if !errors.Is(err, ErrCategoryNotFound) {
		t.Errorf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestClient_GetAllCategories_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetAllCategories(ctx, 1)
	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
}
