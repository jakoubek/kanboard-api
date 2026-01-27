package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetColorList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getColorList" {
			t.Errorf("expected method=getColorList, got %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`{
				"yellow": "Yellow",
				"blue": "Blue",
				"green": "Green",
				"purple": "Purple",
				"red": "Red",
				"orange": "Orange",
				"grey": "Grey",
				"brown": "Brown",
				"deep_orange": "Deep Orange",
				"dark_grey": "Dark Grey",
				"pink": "Pink",
				"teal": "Teal",
				"cyan": "Cyan",
				"lime": "Lime",
				"light_green": "Light Green",
				"amber": "Amber"
			}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	colors, err := client.GetColorList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check that we got some colors
	if len(colors) != 16 {
		t.Errorf("expected 16 colors, got %d", len(colors))
	}

	// Check specific colors
	if colors["yellow"] != "Yellow" {
		t.Errorf("expected yellow='Yellow', got %s", colors["yellow"])
	}
	if colors["blue"] != "Blue" {
		t.Errorf("expected blue='Blue', got %s", colors["blue"])
	}
	if colors["deep_orange"] != "Deep Orange" {
		t.Errorf("expected deep_orange='Deep Orange', got %s", colors["deep_orange"])
	}
}

func TestClient_GetColorList_Empty(t *testing.T) {
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

	colors, err := client.GetColorList(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(colors) != 0 {
		t.Errorf("expected 0 colors, got %d", len(colors))
	}
}
