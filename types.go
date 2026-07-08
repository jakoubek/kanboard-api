package kanboard

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// TaskStatus represents the status of a task.
type TaskStatus int

const (
	// StatusActive represents an open/active task.
	StatusActive TaskStatus = 1
	// StatusInactive represents a closed/inactive task.
	StatusInactive TaskStatus = 0
)

// StringBool is a bool that can be unmarshaled from a string "0" or "1".
type StringBool bool

// UnmarshalJSON implements json.Unmarshaler.
func (b *StringBool) UnmarshalJSON(data []byte) error {
	// Try as string first (most common from Kanboard)
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*b = s == "1" || s == "true"
		return nil
	}

	// Try as number (some Kanboard versions return 0/1)
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*b = n != 0
		return nil
	}

	// Try as raw bool
	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err != nil {
		return err
	}
	*b = StringBool(boolVal)
	return nil
}

// MarshalJSON implements json.Marshaler.
func (b StringBool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`"1"`), nil
	}
	return []byte(`"0"`), nil
}

// StringInt is an int that can be unmarshaled from a JSON string.
type StringInt int

// UnmarshalJSON implements json.Unmarshaler.
func (i *StringInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try as raw int
		var intVal int
		if err := json.Unmarshal(data, &intVal); err != nil {
			return err
		}
		*i = StringInt(intVal)
		return nil
	}
	if s == "" {
		*i = 0
		return nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*i = StringInt(val)
	return nil
}

// StringInt64 is an int64 that can be unmarshaled from a JSON string.
type StringInt64 int64

// UnmarshalJSON implements json.Unmarshaler.
func (i *StringInt64) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try as raw int64
		var intVal int64
		if err := json.Unmarshal(data, &intVal); err != nil {
			return err
		}
		*i = StringInt64(intVal)
		return nil
	}
	if s == "" {
		*i = 0
		return nil
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*i = StringInt64(val)
	return nil
}

// StringFloat is a float64 that can be unmarshaled from a JSON string or number.
type StringFloat float64

// UnmarshalJSON implements json.Unmarshaler.
func (f *StringFloat) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try as raw number
		var num float64
		if err := json.Unmarshal(data, &num); err != nil {
			return err
		}
		*f = StringFloat(num)
		return nil
	}
	if s == "" {
		*f = 0
		return nil
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = StringFloat(val)
	return nil
}

// IntOrFalse is an int that can be unmarshaled from a JSON int or false.
// Kanboard API returns false on failure, int (ID) on success for create operations.
type IntOrFalse int

// UnmarshalJSON implements json.Unmarshaler.
func (i *IntOrFalse) UnmarshalJSON(data []byte) error {
	// Try as int first (success case)
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*i = IntOrFalse(n)
		return nil
	}

	// Try as bool (failure case: false)
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		if b {
			*i = 1 // true shouldn't happen, but handle it
		} else {
			*i = 0
		}
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into IntOrFalse", data)
}

// Project represents a Kanboard project (board).
type Project struct {
	ID                  StringInt  `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	IsActive            StringBool `json:"is_active"`
	Token               string     `json:"token"`
	LastModified        Timestamp  `json:"last_modified"`
	IsPublic            StringBool `json:"is_public"`
	IsPrivate           StringBool `json:"is_private"`
	DefaultSwimlane     string     `json:"default_swimlane"`
	ShowDefaultSwimlane StringBool `json:"show_default_swimlane"`
	Identifier          string     `json:"identifier"`
	StartDate           Timestamp  `json:"start_date"`
	EndDate             Timestamp  `json:"end_date"`
	OwnerID             StringInt  `json:"owner_id"`
	PriorityDefault     StringInt  `json:"priority_default"`
	PriorityStart       StringInt  `json:"priority_start"`
	PriorityEnd         StringInt  `json:"priority_end"`
	Email               string     `json:"email"`
}

// Task represents a Kanboard task.
type Task struct {
	ID                  StringInt   `json:"id"`
	Title               string      `json:"title"`
	Description         string      `json:"description"`
	DateCreation        Timestamp   `json:"date_creation"`
	DateModification    Timestamp   `json:"date_modification"`
	DateCompleted       Timestamp   `json:"date_completed"`
	DateStarted         Timestamp   `json:"date_started"`
	DateDue             Timestamp   `json:"date_due"`
	DateMoved           Timestamp   `json:"date_moved"`
	ColorID             string      `json:"color_id"`
	ProjectID           StringInt   `json:"project_id"`
	ColumnID            StringInt   `json:"column_id"`
	OwnerID             StringInt   `json:"owner_id"`
	CreatorID           StringInt   `json:"creator_id"`
	Position            StringInt   `json:"position"`
	IsActive            StringBool  `json:"is_active"`
	Score               StringInt   `json:"score"`
	TimeEstimated       StringFloat `json:"time_estimated"`
	TimeSpent           StringFloat `json:"time_spent"`
	CategoryID          StringInt   `json:"category_id"`
	SwimlaneID          StringInt   `json:"swimlane_id"`
	Priority            StringInt   `json:"priority"`
	Reference           string      `json:"reference"`
	RecurrenceStatus    StringInt   `json:"recurrence_status"`
	RecurrenceTrigger   StringInt   `json:"recurrence_trigger"`
	RecurrenceFactor    StringInt   `json:"recurrence_factor"`
	RecurrenceTimeframe StringInt   `json:"recurrence_timeframe"`
	RecurrenceBasedate  StringInt   `json:"recurrence_basedate"`
	RecurrenceParent    StringInt   `json:"recurrence_parent"`
	RecurrenceChild     StringInt   `json:"recurrence_child"`
}

// Column represents a Kanboard column.
type Column struct {
	ID          StringInt `json:"id"`
	Title       string    `json:"title"`
	Position    StringInt `json:"position"`
	ProjectID   StringInt `json:"project_id"`
	TaskLimit   StringInt `json:"task_limit"`
	Description string    `json:"description"`
}

// Category represents a Kanboard category.
type Category struct {
	ID        StringInt `json:"id"`
	Name      string    `json:"name"`
	ProjectID StringInt `json:"project_id"`
	ColorID   string    `json:"color_id"`
}

// Comment represents a Kanboard comment.
type Comment struct {
	ID           StringInt `json:"id"`
	TaskID       StringInt `json:"task_id"`
	UserID       StringInt `json:"user_id"`
	DateCreation Timestamp `json:"date_creation"`
	Content      string    `json:"comment"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	AvatarPath   string    `json:"avatar_path"`
}

// TaskLink represents a link between two tasks. Kanboard's getAllTaskLinks returns
// the linked task's ID in the "task_id" field (not "opposite_task_id" or "link_id",
// which don't exist in the real API response).
type TaskLink struct {
	ID             StringInt `json:"id"`
	OppositeTaskID StringInt `json:"task_id"`
	Label          string    `json:"label"`
	Title          string    `json:"title"`
}

// TaskFile represents a file attached to a task.
type TaskFile struct {
	ID           StringInt   `json:"id"`
	Name         string      `json:"name"`
	Path         string      `json:"path"`
	IsImage      StringBool  `json:"is_image"`
	TaskID       StringInt   `json:"task_id"`
	DateCreation Timestamp   `json:"date_creation"`
	UserID       StringInt   `json:"user_id"`
	Size         StringInt64 `json:"size"`
	Username     string      `json:"username"`  // Only returned by getAllTaskFiles
	UserName     string      `json:"user_name"` // Only returned by getAllTaskFiles
}

// Tag represents a Kanboard tag.
type Tag struct {
	ID        StringInt `json:"id"`
	Name      string    `json:"name"`
	ProjectID StringInt `json:"project_id"`
	ColorID   string    `json:"color_id"`
}

// CreateTaskRequest is the request payload for creating a task.
type CreateTaskRequest struct {
	Title               string   `json:"title"`
	ProjectID           int      `json:"project_id"`
	Description         string   `json:"description,omitempty"`
	ColumnID            int      `json:"column_id,omitempty"`
	OwnerID             int      `json:"owner_id,omitempty"`
	CreatorID           int      `json:"creator_id,omitempty"`
	ColorID             string   `json:"color_id,omitempty"`
	CategoryID          int      `json:"category_id,omitempty"`
	DateDue             string   `json:"date_due,omitempty"`
	Score               int      `json:"score,omitempty"`
	SwimlaneID          int      `json:"swimlane_id,omitempty"`
	Priority            int      `json:"priority,omitempty"`
	Reference           string   `json:"reference,omitempty"`
	Tags                []string `json:"tags,omitempty"`
	DateStarted         string   `json:"date_started,omitempty"`
	RecurrenceStatus    int      `json:"recurrence_status,omitempty"`
	RecurrenceTrigger   int      `json:"recurrence_trigger,omitempty"`
	RecurrenceFactor    int      `json:"recurrence_factor,omitempty"`
	RecurrenceTimeframe int      `json:"recurrence_timeframe,omitempty"`
	RecurrenceBasedate  int      `json:"recurrence_basedate,omitempty"`
}

// UpdateTaskRequest is the request payload for updating a task.
// Pointer fields allow distinguishing between "not set" and "set to zero value".
type UpdateTaskRequest struct {
	ID                  int      `json:"id"`
	Title               *string  `json:"title,omitempty"`
	Description         *string  `json:"description,omitempty"`
	ColorID             *string  `json:"color_id,omitempty"`
	OwnerID             *int     `json:"owner_id,omitempty"`
	CategoryID          *int     `json:"category_id,omitempty"`
	DateDue             *string  `json:"date_due,omitempty"`
	Score               *int     `json:"score,omitempty"`
	Priority            *int     `json:"priority,omitempty"`
	Reference           *string  `json:"reference,omitempty"`
	DateStarted         *string  `json:"date_started,omitempty"`
	RecurrenceStatus    *int     `json:"recurrence_status,omitempty"`
	RecurrenceTrigger   *int     `json:"recurrence_trigger,omitempty"`
	RecurrenceFactor    *int     `json:"recurrence_factor,omitempty"`
	RecurrenceTimeframe *int     `json:"recurrence_timeframe,omitempty"`
	RecurrenceBasedate  *int     `json:"recurrence_basedate,omitempty"`
	Tags                []string `json:"tags,omitempty"`
}
