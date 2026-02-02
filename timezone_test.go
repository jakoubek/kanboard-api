package kanboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetTimezone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method != "getTimezone" {
			t.Errorf("expected method=getTimezone, got %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`"Europe/Berlin"`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	tz, err := client.GetTimezone(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tz != "Europe/Berlin" {
		t.Errorf("expected Europe/Berlin, got %s", tz)
	}
}

func TestClient_WithTimezone_ConvertsTaskTimestamps(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		var result json.RawMessage
		switch req.Method {
		case "getTimezone":
			callCount++
			result = json.RawMessage(`"America/New_York"`)
		case "getTask":
			result = json.RawMessage(`{
				"id": "1",
				"title": "Test",
				"description": "",
				"date_creation": 1609459200,
				"date_modification": 1609459200,
				"date_completed": 0,
				"date_started": 0,
				"date_due": 0,
				"date_moved": 0,
				"color_id": "yellow",
				"project_id": "1",
				"column_id": "1",
				"owner_id": "0",
				"creator_id": "1",
				"position": "1",
				"is_active": "1",
				"score": "0",
				"category_id": "0",
				"swimlane_id": "0",
				"priority": "0",
				"reference": "",
				"recurrence_status": "0",
				"recurrence_trigger": "0",
				"recurrence_factor": "0",
				"recurrence_timeframe": "0",
				"recurrence_basedate": "0",
				"recurrence_parent": "0",
				"recurrence_child": "0"
			}`)
		default:
			t.Errorf("unexpected method: %s", req.Method)
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token").WithTimezone()

	task, err := client.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	expected := time.Unix(1609459200, 0).In(loc)
	if !task.DateCreation.Time.Equal(expected) {
		t.Errorf("expected time %v, got %v", expected, task.DateCreation.Time)
	}
	if task.DateCreation.Time.Location().String() != "America/New_York" {
		t.Errorf("expected location America/New_York, got %s", task.DateCreation.Time.Location())
	}

	// Verify getTimezone was called exactly once
	if callCount != 1 {
		t.Errorf("expected getTimezone called once, got %d", callCount)
	}

	// Make a second call — should NOT call getTimezone again
	_, err = client.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected getTimezone still called once, got %d", callCount)
	}
}

func TestClient_WithTimezone_Disabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Method == "getTimezone" {
			t.Error("getTimezone should not be called when timezone is disabled")
		}

		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: json.RawMessage(`{
				"id": "1",
				"title": "Test",
				"description": "",
				"date_creation": 1609459200,
				"date_modification": 0,
				"date_completed": 0,
				"date_started": 0,
				"date_due": 0,
				"date_moved": 0,
				"color_id": "yellow",
				"project_id": "1",
				"column_id": "1",
				"owner_id": "0",
				"creator_id": "1",
				"position": "1",
				"is_active": "1",
				"score": "0",
				"category_id": "0",
				"swimlane_id": "0",
				"priority": "0",
				"reference": "",
				"recurrence_status": "0",
				"recurrence_trigger": "0",
				"recurrence_factor": "0",
				"recurrence_timeframe": "0",
				"recurrence_basedate": "0",
				"recurrence_parent": "0",
				"recurrence_child": "0"
			}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL).WithAPIToken("test-token")

	task, err := client.GetTask(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Without WithTimezone, timestamps stay as unmarshalled (Local from time.Unix)
	if task.DateCreation.Time.Location() != time.Local {
		t.Errorf("expected Local, got %s", task.DateCreation.Time.Location())
	}
}

func TestConvertTimestamps(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	client := &Client{timezone: loc}

	t.Run("struct with Timestamp fields", func(t *testing.T) {
		task := &Task{
			DateCreation:     Timestamp{Time: time.Unix(1609459200, 0)},
			DateModification: Timestamp{Time: time.Unix(1609459200, 0)},
			DateCompleted:    Timestamp{}, // zero — should stay zero
		}
		client.convertTimestamps(task)

		if task.DateCreation.Time.Location().String() != "Asia/Tokyo" {
			t.Errorf("expected Asia/Tokyo, got %s", task.DateCreation.Time.Location())
		}
		if task.DateModification.Time.Location().String() != "Asia/Tokyo" {
			t.Errorf("expected Asia/Tokyo, got %s", task.DateModification.Time.Location())
		}
		if !task.DateCompleted.IsZero() {
			t.Error("zero timestamp should remain zero")
		}
	})

	t.Run("slice of structs", func(t *testing.T) {
		tasks := &[]Task{
			{DateCreation: Timestamp{Time: time.Unix(1609459200, 0)}},
			{DateCreation: Timestamp{Time: time.Unix(1609459200, 0)}},
		}
		client.convertTimestamps(tasks)

		for i, task := range *tasks {
			if task.DateCreation.Time.Location().String() != "Asia/Tokyo" {
				t.Errorf("task[%d]: expected Asia/Tokyo, got %s", i, task.DateCreation.Time.Location())
			}
		}
	})

	t.Run("nil timezone is no-op", func(t *testing.T) {
		noTzClient := &Client{}
		task := &Task{DateCreation: Timestamp{Time: time.Unix(1609459200, 0)}}
		noTzClient.convertTimestamps(task)
		// Should not panic or change anything
		if task.DateCreation.Time.Location() != time.Local {
			t.Errorf("expected Local, got %s", task.DateCreation.Time.Location())
		}
	})
}
