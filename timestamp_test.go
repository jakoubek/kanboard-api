package kanboard

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestamp_UnmarshalJSON_Integer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "positive unix timestamp",
			input:    "1609459200",
			expected: time.Unix(1609459200, 0), // 2021-01-01 00:00:00 UTC
		},
		{
			name:     "zero",
			input:    "0",
			expected: time.Time{},
		},
		{
			name:     "recent timestamp",
			input:    "1704067200",
			expected: time.Unix(1704067200, 0), // 2024-01-01 00:00:00 UTC
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			if err := json.Unmarshal([]byte(tt.input), &ts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ts.Time.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, ts.Time)
			}
		})
	}
}

func TestTimestamp_UnmarshalJSON_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "empty string",
			input:    `""`,
			expected: time.Time{},
		},
		{
			name:     "zero string",
			input:    `"0"`,
			expected: time.Time{},
		},
		{
			name:     "numeric string",
			input:    `"1609459200"`,
			expected: time.Unix(1609459200, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			if err := json.Unmarshal([]byte(tt.input), &ts); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ts.Time.Equal(tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, ts.Time)
			}
		})
	}
}

func TestTimestamp_UnmarshalJSON_Null(t *testing.T) {
	var ts Timestamp
	if err := json.Unmarshal([]byte("null"), &ts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ts.IsZero() {
		t.Errorf("expected zero time for null, got %v", ts.Time)
	}
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		ts       Timestamp
		expected string
	}{
		{
			name:     "zero time",
			ts:       Timestamp{},
			expected: "0",
		},
		{
			name:     "positive timestamp",
			ts:       Timestamp{Time: time.Unix(1609459200, 0)},
			expected: "1609459200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.ts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

func TestTimestamp_RoundTrip(t *testing.T) {
	// Test that marshal/unmarshal produces the same result
	original := Timestamp{Time: time.Unix(1609459200, 0)}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed Timestamp
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !original.Time.Equal(parsed.Time) {
		t.Errorf("round trip failed: %v != %v", original.Time, parsed.Time)
	}
}

func TestTimestamp_InStruct(t *testing.T) {
	// Test Timestamp as part of a struct (simulating API response)
	type Task struct {
		ID           int       `json:"id"`
		DateCreation Timestamp `json:"date_creation"`
		DateDue      Timestamp `json:"date_due"`
	}

	jsonData := `{"id":42,"date_creation":1609459200,"date_due":0}`

	var task Task
	if err := json.Unmarshal([]byte(jsonData), &task); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if task.ID != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if task.DateCreation.IsZero() {
		t.Error("DateCreation should not be zero")
	}
	if !task.DateDue.IsZero() {
		t.Error("DateDue should be zero")
	}
}

func TestTimestamp_InStructWithStringTimestamp(t *testing.T) {
	// Test with string timestamps (Kanboard sometimes returns these)
	type Task struct {
		ID           int       `json:"id"`
		DateCreation Timestamp `json:"date_creation"`
		DateDue      Timestamp `json:"date_due"`
	}

	jsonData := `{"id":42,"date_creation":"1609459200","date_due":""}`

	var task Task
	if err := json.Unmarshal([]byte(jsonData), &task); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if task.ID != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if task.DateCreation.IsZero() {
		t.Error("DateCreation should not be zero")
	}
	if !task.DateDue.IsZero() {
		t.Error("DateDue should be zero")
	}
}

func TestTimestamp_IsZero(t *testing.T) {
	var zero Timestamp
	if !zero.IsZero() {
		t.Error("default Timestamp should be zero")
	}

	nonZero := Timestamp{Time: time.Unix(1609459200, 0)}
	if nonZero.IsZero() {
		t.Error("non-zero Timestamp should not be zero")
	}
}
