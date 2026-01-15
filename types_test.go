package kanboard

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTaskStatus_Constants(t *testing.T) {
	if StatusActive != 1 {
		t.Errorf("expected StatusActive=1, got %d", StatusActive)
	}
	if StatusInactive != 0 {
		t.Errorf("expected StatusInactive=0, got %d", StatusInactive)
	}
}

func TestStringBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"string 1", `"1"`, true},
		{"string 0", `"0"`, false},
		{"string true", `"true"`, true},
		{"string false", `"false"`, false},
		{"number 1", `1`, true},
		{"number 0", `0`, false},
		{"number non-zero", `42`, true},
		{"bool true", `true`, true},
		{"bool false", `false`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b StringBool
			if err := json.Unmarshal([]byte(tt.input), &b); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if bool(b) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, b)
			}
		})
	}
}

func TestStringInt_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"string number", `"42"`, 42},
		{"string zero", `"0"`, 0},
		{"int number", `42`, 42},
		{"int zero", `0`, 0},
		{"empty string", `""`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i StringInt
			if err := json.Unmarshal([]byte(tt.input), &i); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if int(i) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, i)
			}
		})
	}
}

func TestStringInt64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"string number", `"1048576"`, 1048576},
		{"string zero", `"0"`, 0},
		{"int64 number", `1048576`, 1048576},
		{"int64 zero", `0`, 0},
		{"empty string", `""`, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i StringInt64
			if err := json.Unmarshal([]byte(tt.input), &i); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if int64(i) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, i)
			}
		})
	}
}

func TestIntOrFalse_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"int value", `42`, 42},
		{"int zero", `0`, 0},
		{"false", `false`, 0},
		{"true", `true`, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var i IntOrFalse
			if err := json.Unmarshal([]byte(tt.input), &i); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if int(i) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, i)
			}
		})
	}
}

func TestProject_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "1",
		"name": "Test Project",
		"description": "A test project",
		"is_active": "1",
		"token": "abc123",
		"last_modified": 1609459200,
		"is_public": "0",
		"is_private": "1",
		"owner_id": "42",
		"priority_default": "2"
	}`

	var project Project
	if err := json.Unmarshal([]byte(jsonData), &project); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(project.ID) != 1 {
		t.Errorf("expected ID=1, got %d", project.ID)
	}
	if project.Name != "Test Project" {
		t.Errorf("expected Name='Test Project', got %s", project.Name)
	}
	if !bool(project.IsActive) {
		t.Error("expected IsActive=true")
	}
	if bool(project.IsPublic) {
		t.Error("expected IsPublic=false")
	}
	if !bool(project.IsPrivate) {
		t.Error("expected IsPrivate=true")
	}
	if int(project.OwnerID) != 42 {
		t.Errorf("expected OwnerID=42, got %d", project.OwnerID)
	}
}

func TestProject_UnmarshalJSON_NumericBool(t *testing.T) {
	// Some Kanboard versions return numeric booleans instead of strings
	jsonData := `{
		"id": 1,
		"name": "Test Project",
		"description": "A test project",
		"is_active": 1,
		"token": "abc123",
		"last_modified": 1609459200,
		"is_public": 0,
		"is_private": 1,
		"owner_id": 42,
		"priority_default": 2
	}`

	var project Project
	if err := json.Unmarshal([]byte(jsonData), &project); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(project.ID) != 1 {
		t.Errorf("expected ID=1, got %d", project.ID)
	}
	if !bool(project.IsActive) {
		t.Error("expected IsActive=true")
	}
	if bool(project.IsPublic) {
		t.Error("expected IsPublic=false")
	}
	if !bool(project.IsPrivate) {
		t.Error("expected IsPrivate=true")
	}
}

func TestTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "42",
		"title": "Test Task",
		"description": "Task description",
		"date_creation": 1609459200,
		"date_due": 0,
		"color_id": "yellow",
		"project_id": "1",
		"column_id": "2",
		"owner_id": "5",
		"is_active": "1",
		"priority": "2",
		"category_id": "3"
	}`

	var task Task
	if err := json.Unmarshal([]byte(jsonData), &task); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(task.ID) != 42 {
		t.Errorf("expected ID=42, got %d", task.ID)
	}
	if task.Title != "Test Task" {
		t.Errorf("expected Title='Test Task', got %s", task.Title)
	}
	if int(task.ProjectID) != 1 {
		t.Errorf("expected ProjectID=1, got %d", task.ProjectID)
	}
	if int(task.ColumnID) != 2 {
		t.Errorf("expected ColumnID=2, got %d", task.ColumnID)
	}
	if !bool(task.IsActive) {
		t.Error("expected IsActive=true")
	}
	if task.DateCreation.IsZero() {
		t.Error("expected DateCreation to be set")
	}
	if !task.DateDue.IsZero() {
		t.Error("expected DateDue to be zero")
	}
}

func TestColumn_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "1",
		"title": "Backlog",
		"position": "1",
		"project_id": "5",
		"task_limit": "10",
		"description": "Tasks to be done"
	}`

	var column Column
	if err := json.Unmarshal([]byte(jsonData), &column); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(column.ID) != 1 {
		t.Errorf("expected ID=1, got %d", column.ID)
	}
	if column.Title != "Backlog" {
		t.Errorf("expected Title='Backlog', got %s", column.Title)
	}
	if int(column.Position) != 1 {
		t.Errorf("expected Position=1, got %d", column.Position)
	}
	if int(column.TaskLimit) != 10 {
		t.Errorf("expected TaskLimit=10, got %d", column.TaskLimit)
	}
}

func TestCategory_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "3",
		"name": "Bug",
		"project_id": "1",
		"color_id": "red"
	}`

	var category Category
	if err := json.Unmarshal([]byte(jsonData), &category); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(category.ID) != 3 {
		t.Errorf("expected ID=3, got %d", category.ID)
	}
	if category.Name != "Bug" {
		t.Errorf("expected Name='Bug', got %s", category.Name)
	}
	if category.ColorID != "red" {
		t.Errorf("expected ColorID='red', got %s", category.ColorID)
	}
}

func TestComment_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "10",
		"task_id": "42",
		"user_id": "5",
		"date_creation": 1609459200,
		"comment": "This is a comment",
		"username": "admin",
		"name": "Admin User",
		"email": "admin@example.com"
	}`

	var comment Comment
	if err := json.Unmarshal([]byte(jsonData), &comment); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(comment.ID) != 10 {
		t.Errorf("expected ID=10, got %d", comment.ID)
	}
	if int(comment.TaskID) != 42 {
		t.Errorf("expected TaskID=42, got %d", comment.TaskID)
	}
	if comment.Content != "This is a comment" {
		t.Errorf("expected Content='This is a comment', got %s", comment.Content)
	}
	if comment.Username != "admin" {
		t.Errorf("expected Username='admin', got %s", comment.Username)
	}
}

func TestTaskLink_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "1",
		"link_id": "2",
		"task_id": "42",
		"opposite_task_id": "43",
		"label": "blocks",
		"title": "Related Task"
	}`

	var link TaskLink
	if err := json.Unmarshal([]byte(jsonData), &link); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(link.ID) != 1 {
		t.Errorf("expected ID=1, got %d", link.ID)
	}
	if int(link.TaskID) != 42 {
		t.Errorf("expected TaskID=42, got %d", link.TaskID)
	}
	if int(link.OppositeTaskID) != 43 {
		t.Errorf("expected OppositeTaskID=43, got %d", link.OppositeTaskID)
	}
	if link.Label != "blocks" {
		t.Errorf("expected Label='blocks', got %s", link.Label)
	}
}

func TestTaskFile_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "1",
		"name": "document.pdf",
		"path": "/uploads/document.pdf",
		"is_image": "0",
		"task_id": "42",
		"date_creation": 1609459200,
		"user_id": "5",
		"size": "1048576"
	}`

	var file TaskFile
	if err := json.Unmarshal([]byte(jsonData), &file); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(file.ID) != 1 {
		t.Errorf("expected ID=1, got %d", file.ID)
	}
	if file.Name != "document.pdf" {
		t.Errorf("expected Name='document.pdf', got %s", file.Name)
	}
	if bool(file.IsImage) {
		t.Error("expected IsImage=false")
	}
	if int64(file.Size) != 1048576 {
		t.Errorf("expected Size=1048576, got %d", file.Size)
	}
}

func TestTag_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"id": "5",
		"name": "urgent",
		"project_id": "1",
		"color_id": "red"
	}`

	var tag Tag
	if err := json.Unmarshal([]byte(jsonData), &tag); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if int(tag.ID) != 5 {
		t.Errorf("expected ID=5, got %d", tag.ID)
	}
	if tag.Name != "urgent" {
		t.Errorf("expected Name='urgent', got %s", tag.Name)
	}
	if tag.ColorID != "red" {
		t.Errorf("expected ColorID='red', got %s", tag.ColorID)
	}
}

func TestCreateTaskRequest_MarshalJSON(t *testing.T) {
	req := CreateTaskRequest{
		Title:       "New Task",
		ProjectID:   1,
		Description: "Task description",
		ColumnID:    2,
		Priority:    3,
		Tags:        []string{"urgent", "backend"},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if unmarshaled["title"] != "New Task" {
		t.Errorf("expected title='New Task', got %v", unmarshaled["title"])
	}
	if unmarshaled["project_id"].(float64) != 1 {
		t.Errorf("expected project_id=1, got %v", unmarshaled["project_id"])
	}
	tags := unmarshaled["tags"].([]any)
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}
}

func TestCreateTaskRequest_OmitEmpty(t *testing.T) {
	req := CreateTaskRequest{
		Title:     "Minimal Task",
		ProjectID: 1,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// Should not have omitempty fields
	if _, exists := unmarshaled["description"]; exists {
		t.Error("description should be omitted when empty")
	}
	if _, exists := unmarshaled["color_id"]; exists {
		t.Error("color_id should be omitted when empty")
	}
}

func TestUpdateTaskRequest_MarshalJSON(t *testing.T) {
	title := "Updated Title"
	priority := 5

	req := UpdateTaskRequest{
		ID:       42,
		Title:    &title,
		Priority: &priority,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if unmarshaled["id"].(float64) != 42 {
		t.Errorf("expected id=42, got %v", unmarshaled["id"])
	}
	if unmarshaled["title"] != "Updated Title" {
		t.Errorf("expected title='Updated Title', got %v", unmarshaled["title"])
	}
	if unmarshaled["priority"].(float64) != 5 {
		t.Errorf("expected priority=5, got %v", unmarshaled["priority"])
	}

	// Should not have fields that weren't set
	if _, exists := unmarshaled["description"]; exists {
		t.Error("description should be omitted when nil")
	}
}

func TestUpdateTaskRequest_ZeroValueVsNil(t *testing.T) {
	// Test that we can distinguish between "not set" and "set to zero"
	zero := 0
	emptyString := ""

	req := UpdateTaskRequest{
		ID:       42,
		Priority: &zero,        // Explicitly set to 0
		ColorID:  &emptyString, // Explicitly set to empty string
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// These should be present even though they're zero values
	if _, exists := unmarshaled["priority"]; !exists {
		t.Error("priority should be present when explicitly set to 0")
	}
	if _, exists := unmarshaled["color_id"]; !exists {
		t.Error("color_id should be present when explicitly set to empty string")
	}
}

func TestTask_TimestampFields(t *testing.T) {
	jsonData := `{
		"id": "1",
		"title": "Test",
		"description": "",
		"date_creation": 1609459200,
		"date_modification": 1609545600,
		"date_completed": 0,
		"date_due": 1610064000,
		"color_id": "blue",
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
	}`

	var task Task
	if err := json.Unmarshal([]byte(jsonData), &task); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// Check that creation date is parsed correctly
	expectedCreation := time.Unix(1609459200, 0)
	if !task.DateCreation.Time.Equal(expectedCreation) {
		t.Errorf("expected DateCreation=%v, got %v", expectedCreation, task.DateCreation.Time)
	}

	// Check that completed date is zero
	if !task.DateCompleted.IsZero() {
		t.Error("expected DateCompleted to be zero")
	}

	// Check that due date is parsed
	if task.DateDue.IsZero() {
		t.Error("expected DateDue to be set")
	}
}
