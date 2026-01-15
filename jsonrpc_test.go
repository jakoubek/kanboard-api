package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestJSONRPCRequest_Marshal(t *testing.T) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "getTask",
		ID:      1,
		Params:  map[string]int{"task_id": 42},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if unmarshaled["jsonrpc"] != "2.0" {
		t.Errorf("expected jsonrpc=2.0, got %v", unmarshaled["jsonrpc"])
	}
	if unmarshaled["method"] != "getTask" {
		t.Errorf("expected method=getTask, got %v", unmarshaled["method"])
	}
	if unmarshaled["id"].(float64) != 1 {
		t.Errorf("expected id=1, got %v", unmarshaled["id"])
	}
}

func TestJSONRPCRequest_MarshalWithoutParams(t *testing.T) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "getAllProjects",
		ID:      1,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal request: %v", err)
	}

	if _, exists := unmarshaled["params"]; exists {
		t.Error("params should be omitted when nil")
	}
}

func TestJSONRPCResponse_Unmarshal(t *testing.T) {
	data := `{"jsonrpc":"2.0","id":1,"result":{"id":42,"title":"Test Task"}}`

	var resp JSONRPCResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc=2.0, got %v", resp.JSONRPC)
	}
	if resp.ID != 1 {
		t.Errorf("expected id=1, got %v", resp.ID)
	}
	if resp.Error != nil {
		t.Error("expected no error")
	}
	if resp.Result == nil {
		t.Error("expected result to be present")
	}
}

func TestJSONRPCResponse_UnmarshalError(t *testing.T) {
	data := `{"jsonrpc":"2.0","id":1,"error":{"code":-32600,"message":"Invalid Request"}}`

	var resp JSONRPCResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error to be present")
	}
	if resp.Error.Code != -32600 {
		t.Errorf("expected error code=-32600, got %v", resp.Error.Code)
	}
	if resp.Error.Message != "Invalid Request" {
		t.Errorf("expected error message='Invalid Request', got %v", resp.Error.Message)
	}
}

func TestJSONRPCError_Error(t *testing.T) {
	err := &JSONRPCError{
		Code:    -32600,
		Message: "Invalid Request",
	}

	expected := "JSON-RPC error (code -32600): Invalid Request"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestNextRequestID_Increments(t *testing.T) {
	// Get the current counter value
	initial := nextRequestID()

	// Verify increments
	for i := int64(1); i <= 5; i++ {
		got := nextRequestID()
		expected := initial + i
		if got != expected {
			t.Errorf("expected %d, got %d", expected, got)
		}
	}
}

func TestNextRequestID_ThreadSafe(t *testing.T) {
	const goroutines = 100
	const iterations = 100

	var wg sync.WaitGroup
	ids := make(chan int64, goroutines*iterations)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ids <- nextRequestID()
			}
		}()
	}

	wg.Wait()
	close(ids)

	// Collect all IDs and check for uniqueness
	seen := make(map[int64]bool)
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate request ID: %d", id)
		}
		seen[id] = true
	}

	if len(seen) != goroutines*iterations {
		t.Errorf("expected %d unique IDs, got %d", goroutines*iterations, len(seen))
	}
}

func TestClient_Call_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/jsonrpc.php" {
			t.Errorf("expected /jsonrpc.php, got %s", r.URL.Path)
		}

		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc=2.0, got %s", req.JSONRPC)
		}
		if req.Method != "getTask" {
			t.Errorf("expected method=getTask, got %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"id":42,"title":"Test Task"}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	var result struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	}

	err := client.call(context.Background(), "getTask", map[string]int{"task_id": 42}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != 42 {
		t.Errorf("expected id=42, got %d", result.ID)
	}
	if result.Title != "Test Task" {
		t.Errorf("expected title='Test Task', got %s", result.Title)
	}
}

func TestClient_Call_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &JSONRPCError{
				Code:    -32600,
				Message: "Invalid Request",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	var result interface{}
	err := client.call(context.Background(), "invalidMethod", nil, &result)

	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Code != -32600 {
		t.Errorf("expected code=-32600, got %d", apiErr.Code)
	}
}

func TestClient_Call_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("invalid-token")

	var result interface{}
	err := client.call(context.Background(), "getTask", nil, &result)

	if err != ErrUnauthorized {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestClient_Call_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	var result interface{}
	err := client.call(context.Background(), "getTask", nil, &result)

	if err != ErrForbidden {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestClient_Call_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		select {}
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result interface{}
	err := client.call(ctx, "getTask", nil, &result)

	if err == nil {
		t.Fatal("expected error due to canceled context")
	}
}

func TestClient_Call_SubdirectoryInstallation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/kanboard/jsonrpc.php" {
			t.Errorf("expected /kanboard/jsonrpc.php, got %s", r.URL.Path)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Client with subdirectory path
	client := NewClient(server.URL + "/kanboard").WithAPIToken("test-token")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Call_TrailingSlashHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/jsonrpc.php" {
			t.Errorf("expected /jsonrpc.php, got %s", r.URL.Path)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Client with trailing slash
	client := NewClient(server.URL + "/").WithAPIToken("test-token")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Call_AuthHeaderSent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected Basic Auth header")
		}
		if username != "jsonrpc" {
			t.Errorf("expected username=jsonrpc, got %s", username)
		}
		if password != "test-token" {
			t.Errorf("expected password=test-token, got %s", password)
		}

		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`true`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	var result bool
	err := client.call(context.Background(), "getVersion", nil, &result)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
