package kanboard

import "time"

// TaskUpdateParams is a fluent builder for task update configuration.
// Only set fields are included in the update request.
type TaskUpdateParams struct {
	title       *string
	description *string
	colorID     *string
	ownerID     *int
	categoryID  *int
	priority    *int
	score       *int
	dueDate     *int64
	startDate   *int64
	reference   *string
	tags        []string
	tagsSet     bool // tracks whether tags were explicitly set (even to empty)
}

// NewTaskUpdate creates a new TaskUpdateParams.
func NewTaskUpdate() *TaskUpdateParams {
	return &TaskUpdateParams{}
}

// SetTitle sets the task title.
func (p *TaskUpdateParams) SetTitle(title string) *TaskUpdateParams {
	p.title = &title
	return p
}

// SetDescription sets the task description.
func (p *TaskUpdateParams) SetDescription(desc string) *TaskUpdateParams {
	p.description = &desc
	return p
}

// SetColor sets the color ID for the task.
func (p *TaskUpdateParams) SetColor(colorID string) *TaskUpdateParams {
	p.colorID = &colorID
	return p
}

// SetOwner sets the owner (assignee) ID for the task.
func (p *TaskUpdateParams) SetOwner(ownerID int) *TaskUpdateParams {
	p.ownerID = &ownerID
	return p
}

// SetCategory sets the category ID for the task.
func (p *TaskUpdateParams) SetCategory(categoryID int) *TaskUpdateParams {
	p.categoryID = &categoryID
	return p
}

// SetPriority sets the priority for the task.
func (p *TaskUpdateParams) SetPriority(priority int) *TaskUpdateParams {
	p.priority = &priority
	return p
}

// SetScore sets the complexity score for the task.
func (p *TaskUpdateParams) SetScore(score int) *TaskUpdateParams {
	p.score = &score
	return p
}

// SetDueDate sets the due date for the task.
func (p *TaskUpdateParams) SetDueDate(date time.Time) *TaskUpdateParams {
	ts := date.Unix()
	p.dueDate = &ts
	return p
}

// SetStartDate sets the start date for the task.
func (p *TaskUpdateParams) SetStartDate(date time.Time) *TaskUpdateParams {
	ts := date.Unix()
	p.startDate = &ts
	return p
}

// SetReference sets the external reference for the task.
func (p *TaskUpdateParams) SetReference(ref string) *TaskUpdateParams {
	p.reference = &ref
	return p
}

// SetTags sets the tags for the task.
// This replaces all existing tags on the task.
// Call with no arguments to clear all tags.
func (p *TaskUpdateParams) SetTags(tags ...string) *TaskUpdateParams {
	if tags == nil {
		p.tags = []string{}
	} else {
		p.tags = tags
	}
	p.tagsSet = true
	return p
}

// ClearDueDate clears the due date from the task.
func (p *TaskUpdateParams) ClearDueDate() *TaskUpdateParams {
	zero := int64(0)
	p.dueDate = &zero
	return p
}

// ClearStartDate clears the start date from the task.
func (p *TaskUpdateParams) ClearStartDate() *TaskUpdateParams {
	zero := int64(0)
	p.startDate = &zero
	return p
}

// ClearOwner removes the owner from the task.
func (p *TaskUpdateParams) ClearOwner() *TaskUpdateParams {
	zero := 0
	p.ownerID = &zero
	return p
}

// ClearCategory removes the category from the task.
func (p *TaskUpdateParams) ClearCategory() *TaskUpdateParams {
	zero := 0
	p.categoryID = &zero
	return p
}

// toUpdateTaskRequest converts TaskUpdateParams to an UpdateTaskRequest.
// The taskID is required and must be provided by the caller.
func (p *TaskUpdateParams) toUpdateTaskRequest(taskID int) UpdateTaskRequest {
	req := UpdateTaskRequest{
		ID: taskID,
	}

	if p.title != nil {
		req.Title = p.title
	}
	if p.description != nil {
		req.Description = p.description
	}
	if p.colorID != nil {
		req.ColorID = p.colorID
	}
	if p.ownerID != nil {
		req.OwnerID = p.ownerID
	}
	if p.categoryID != nil {
		req.CategoryID = p.categoryID
	}
	if p.priority != nil {
		req.Priority = p.priority
	}
	if p.score != nil {
		req.Score = p.score
	}
	if p.dueDate != nil {
		req.DateDue = p.dueDate
	}
	if p.startDate != nil {
		req.DateStarted = p.startDate
	}
	if p.reference != nil {
		req.Reference = p.reference
	}
	if p.tagsSet {
		req.Tags = p.tags
	}

	return req
}
