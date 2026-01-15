package kanboard

import (
	"encoding/json"
	"fmt"
	"time"
)

// Timestamp wraps time.Time and handles Kanboard's Unix timestamp JSON format.
// Kanboard returns timestamps as Unix integers, with 0 or empty strings for null values.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler.
// Supports Unix timestamps as integers, empty strings, "0" strings, and zero values.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as integer (Unix timestamp)
	var unix int64
	if err := json.Unmarshal(data, &unix); err == nil {
		if unix == 0 {
			t.Time = time.Time{}
		} else {
			t.Time = time.Unix(unix, 0)
		}
		return nil
	}

	// Try to unmarshal as string (empty or "0")
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" || str == "0" {
			t.Time = time.Time{}
			return nil
		}
		// Try to parse as numeric string
		var unix int64
		if _, err := fmt.Sscanf(str, "%d", &unix); err == nil {
			if unix == 0 {
				t.Time = time.Time{}
			} else {
				t.Time = time.Unix(unix, 0)
			}
			return nil
		}
	}

	// Handle null
	if string(data) == "null" {
		t.Time = time.Time{}
		return nil
	}

	return fmt.Errorf("cannot unmarshal timestamp: %s", string(data))
}

// MarshalJSON implements json.Marshaler.
// Returns 0 for zero time, otherwise returns Unix timestamp.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("0"), nil
	}
	return json.Marshal(t.Unix())
}
