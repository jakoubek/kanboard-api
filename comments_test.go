package kanboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetAllComments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllComments" {
			t.Errorf("expected method=getAllComments, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "1", "task_id": "42", "user_id": "1", "comment": "First comment", "username": "admin"},
				{"id": "2", "task_id": "42", "user_id": "2", "comment": "Second comment", "username": "user"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	comments, err := client.GetAllComments(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
	if comments[0].Content != "First comment" {
		t.Errorf("expected first comment='First comment', got %s", comments[0].Content)
	}
	if comments[0].Username != "admin" {
		t.Errorf("expected username='admin', got %s", comments[0].Username)
	}
}

func TestClient_GetAllComments_Empty(t *testing.T) {
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

	comments, err := client.GetAllComments(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}

func TestClient_GetComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getComment" {
			t.Errorf("expected method=getComment, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["comment_id"].(float64) != 5 {
			t.Errorf("expected comment_id=5, got %v", params["comment_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id": "5", "task_id": "42", "user_id": "1", "comment": "Test comment", "username": "admin"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	comment, err := client.GetComment(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(comment.ID) != 5 {
		t.Errorf("expected ID=5, got %d", comment.ID)
	}
	if comment.Content != "Test comment" {
		t.Errorf("expected content='Test comment', got %s", comment.Content)
	}
}

func TestClient_GetComment_NotFound(t *testing.T) {
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

	_, err := client.GetComment(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}

	if !errors.Is(err, ErrCommentNotFound) {
		t.Errorf("expected ErrCommentNotFound, got %v", err)
	}
}

func TestClient_CreateComment(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			// First call: createComment
			if req.Method != "createComment" {
				t.Errorf("expected method=createComment, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			if params["task_id"].(float64) != 42 {
				t.Errorf("expected task_id=42, got %v", params["task_id"])
			}
			if params["user_id"].(float64) != 1 {
				t.Errorf("expected user_id=1, got %v", params["user_id"])
			}
			if params["content"] != "New comment" {
				t.Errorf("expected content='New comment', got %v", params["content"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`10`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: getComment to fetch created comment
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "10", "task_id": "42", "user_id": "1", "comment": "New comment", "username": "admin"}`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	comment, err := client.CreateComment(context.Background(), 42, 1, "New comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(comment.ID) != 10 {
		t.Errorf("expected ID=10, got %d", comment.ID)
	}
	if comment.Content != "New comment" {
		t.Errorf("expected content='New comment', got %s", comment.Content)
	}
}

func TestClient_CreateComment_Failure(t *testing.T) {
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

	_, err := client.CreateComment(context.Background(), 42, 1, "New comment")
	if err == nil {
		t.Fatal("expected error for failed comment creation")
	}
}

func TestClient_UpdateComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "updateComment" {
			t.Errorf("expected method=updateComment, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["id"].(float64) != 5 {
			t.Errorf("expected id=5, got %v", params["id"])
		}
		if params["content"] != "Updated content" {
			t.Errorf("expected content='Updated content', got %v", params["content"])
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

	err := client.UpdateComment(context.Background(), 5, "Updated content")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_UpdateComment_Failure(t *testing.T) {
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

	err := client.UpdateComment(context.Background(), 5, "Updated content")
	if err == nil {
		t.Fatal("expected error for failed update")
	}
}

func TestClient_RemoveComment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "removeComment" {
			t.Errorf("expected method=removeComment, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["comment_id"].(float64) != 5 {
			t.Errorf("expected comment_id=5, got %v", params["comment_id"])
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

	err := client.RemoveComment(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_RemoveComment_Failure(t *testing.T) {
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

	err := client.RemoveComment(context.Background(), 5)
	if err == nil {
		t.Fatal("expected error for failed delete")
	}
}

func TestTaskScope_GetComments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllComments" {
			t.Errorf("expected method=getAllComments, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "task_id": "42", "comment": "Test", "username": "admin"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	comments, err := client.Task(42).GetComments(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
}

func TestTaskScope_AddComment(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		callCount++
		if callCount == 1 {
			// First call: createComment
			if req.Method != "createComment" {
				t.Errorf("expected method=createComment, got %s", req.Method)
			}

			params := req.Params.(map[string]any)
			if params["task_id"].(float64) != 42 {
				t.Errorf("expected task_id=42, got %v", params["task_id"])
			}

			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`10`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: getComment
			resp := JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  json.RawMessage(`{"id": "10", "task_id": "42", "comment": "Added via scope", "username": "admin"}`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	comment, err := client.Task(42).AddComment(context.Background(), 1, "Added via scope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if int(comment.ID) != 10 {
		t.Errorf("expected ID=10, got %d", comment.ID)
	}
}
