package kanboard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetAllTaskFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllTaskFiles" {
			t.Errorf("expected method=getAllTaskFiles, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`[
				{"id": "1", "name": "document.pdf", "path": "files/1", "is_image": "0", "task_id": "42", "user_id": "1", "size": "1024"},
				{"id": "2", "name": "screenshot.png", "path": "files/2", "is_image": "1", "task_id": "42", "user_id": "1", "size": "2048"}
			]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	files, err := client.GetAllTaskFiles(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
	if files[0].Name != "document.pdf" {
		t.Errorf("expected name='document.pdf', got %s", files[0].Name)
	}
	if bool(files[1].IsImage) != true {
		t.Error("expected second file to be an image")
	}
}

func TestClient_GetAllTaskFiles_Empty(t *testing.T) {
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

	files, err := client.GetAllTaskFiles(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestClient_CreateTaskFile(t *testing.T) {
	testContent := []byte("test file content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "createTaskFile" {
			t.Errorf("expected method=createTaskFile, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["project_id"].(float64) != 1 {
			t.Errorf("expected project_id=1, got %v", params["project_id"])
		}
		if params["task_id"].(float64) != 42 {
			t.Errorf("expected task_id=42, got %v", params["task_id"])
		}
		if params["filename"] != "test.txt" {
			t.Errorf("expected filename='test.txt', got %v", params["filename"])
		}

		// Verify base64 encoded content
		expectedBlob := base64.StdEncoding.EncodeToString(testContent)
		if params["blob"] != expectedBlob {
			t.Errorf("expected blob to be base64 encoded")
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

	fileID, err := client.CreateTaskFile(context.Background(), 1, 42, "test.txt", testContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fileID != 10 {
		t.Errorf("expected fileID=10, got %d", fileID)
	}
}

func TestClient_CreateTaskFile_Failure(t *testing.T) {
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

	_, err := client.CreateTaskFile(context.Background(), 1, 42, "test.txt", []byte("content"))
	if err == nil {
		t.Fatal("expected error for failed file upload")
	}
}

func TestClient_DownloadTaskFile(t *testing.T) {
	testContent := []byte("downloaded content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "downloadTaskFile" {
			t.Errorf("expected method=downloadTaskFile, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["file_id"].(float64) != 5 {
			t.Errorf("expected file_id=5, got %v", params["file_id"])
		}

		// Return base64 encoded content
		encoded := base64.StdEncoding.EncodeToString(testContent)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`"` + encoded + `"`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	content, err := client.DownloadTaskFile(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("expected content='%s', got '%s'", testContent, content)
	}
}

func TestClient_RemoveTaskFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "removeTaskFile" {
			t.Errorf("expected method=removeTaskFile, got %s", req.Method)
		}

		params := req.Params.(map[string]any)
		if params["file_id"].(float64) != 5 {
			t.Errorf("expected file_id=5, got %v", params["file_id"])
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

	err := client.RemoveTaskFile(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_RemoveTaskFile_Failure(t *testing.T) {
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

	err := client.RemoveTaskFile(context.Background(), 5)
	if err == nil {
		t.Fatal("expected error for failed delete")
	}
}

func TestTaskScope_GetFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getAllTaskFiles" {
			t.Errorf("expected method=getAllTaskFiles, got %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`[{"id": "1", "name": "file.txt", "task_id": "42", "is_image": "0"}]`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	files, err := client.Task(42).GetFiles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

func TestTaskScope_UploadFile(t *testing.T) {
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
				Result:  json.RawMessage(`{"id": "42", "project_id": "1", "is_active": "1"}`),
			}
			json.NewEncoder(w).Encode(resp)
		} else {
			// Second call: createTaskFile
			if req.Method != "createTaskFile" {
				t.Errorf("expected method=createTaskFile, got %s", req.Method)
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
				Result:  json.RawMessage(`10`),
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	fileID, err := client.Task(42).UploadFile(context.Background(), "test.txt", []byte("content"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fileID != 10 {
		t.Errorf("expected fileID=10, got %d", fileID)
	}
}
