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

func TestClient_CreateCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "createCategory" {
			t.Errorf("expected method=createCategory, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["name"].(string) != "Bug" {
			t.Errorf("expected name=Bug, got %v", params["name"])
		}
		if params["color_id"].(string) != "red" {
			t.Errorf("expected color_id=red, got %v", params["color_id"])
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

	id, err := client.CreateCategory(context.Background(), 1, "Bug", "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id=42, got %d", id)
	}
}

func TestClient_CreateCategory_NoColor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		params := req.Params.(map[string]any)
		if _, ok := params["color_id"]; ok {
			t.Error("expected color_id to be absent")
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

	id, err := client.CreateCategory(context.Background(), 1, "Bug", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 10 {
		t.Errorf("expected id=10, got %d", id)
	}
}

func TestClient_CreateCategory_Failure(t *testing.T) {
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

	_, err := client.CreateCategory(context.Background(), 1, "Bug", "")
	if err == nil {
		t.Fatal("expected error on failure")
	}
}

func TestClient_UpdateCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "updateCategory" {
			t.Errorf("expected method=updateCategory, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["id"].(float64) != 5 {
			t.Errorf("expected id=5, got %v", params["id"])
		}
		if params["name"].(string) != "Feature" {
			t.Errorf("expected name=Feature, got %v", params["name"])
		}
		if params["color_id"].(string) != "blue" {
			t.Errorf("expected color_id=blue, got %v", params["color_id"])
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

	err := client.UpdateCategory(context.Background(), 5, "Feature", "blue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_UpdateCategory_Failure(t *testing.T) {
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

	err := client.UpdateCategory(context.Background(), 5, "Feature", "")
	if err == nil {
		t.Fatal("expected error on failure")
	}
}

func TestClient_RemoveCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "removeCategory" {
			t.Errorf("expected method=removeCategory, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["category_id"].(float64) != 3 {
			t.Errorf("expected category_id=3, got %v", params["category_id"])
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

	err := client.RemoveCategory(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_RemoveCategory_Failure(t *testing.T) {
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

	err := client.RemoveCategory(context.Background(), 3)
	if err == nil {
		t.Fatal("expected error on failure")
	}
}

func TestClient_GetCategoryByName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

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

	cat, err := client.GetCategoryByName(context.Background(), 1, "Feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cat.Name != "Feature" {
		t.Errorf("expected name=Feature, got %s", cat.Name)
	}
	if int(cat.ID) != 2 {
		t.Errorf("expected id=2, got %d", cat.ID)
	}
}

func TestClient_GetCategoryByName_NotFound(t *testing.T) {
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

	_, err := client.GetCategoryByName(context.Background(), 1, "Nonexistent")
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
