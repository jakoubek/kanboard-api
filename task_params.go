package kanboard

import "time"

// TaskParams is a fluent builder for task creation configuration.
// It is a pure configuration object with no I/O.
type TaskParams struct {
	title       string
	description *string
	columnID    *int
	categoryID  *int
	ownerID     *int
	creatorID   *int
	colorID     *string
	priority    *int
	score       *int
	dueDate     *int64
	startDate   *int64
	swimlaneID  *int
	reference   *string
	tags        []string
}

// NewTask creates a new TaskParams with the given title.
func NewTask(title string) *TaskParams {
	return &TaskParams{
		title: title,
	}
}

// WithDescription sets the task description.
func (p *TaskParams) WithDescription(desc string) *TaskParams {
	p.description = &desc
	return p
}

// InColumn sets the column ID for the task.
func (p *TaskParams) InColumn(columnID int) *TaskParams {
	p.columnID = &columnID
	return p
}

// WithCategory sets the category ID for the task.
func (p *TaskParams) WithCategory(categoryID int) *TaskParams {
	p.categoryID = &categoryID
	return p
}

// WithOwner sets the owner (assignee) ID for the task.
func (p *TaskParams) WithOwner(ownerID int) *TaskParams {
	p.ownerID = &ownerID
	return p
}

// WithCreator sets the creator ID for the task.
func (p *TaskParams) WithCreator(creatorID int) *TaskParams {
	p.creatorID = &creatorID
	return p
}

// WithColor sets the color ID for the task.
func (p *TaskParams) WithColor(colorID string) *TaskParams {
	p.colorID = &colorID
	return p
}

// WithPriority sets the priority for the task.
func (p *TaskParams) WithPriority(priority int) *TaskParams {
	p.priority = &priority
	return p
}

// WithScore sets the complexity score for the task.
func (p *TaskParams) WithScore(score int) *TaskParams {
	p.score = &score
	return p
}

// WithDueDate sets the due date for the task.
func (p *TaskParams) WithDueDate(date time.Time) *TaskParams {
	ts := date.Unix()
	p.dueDate = &ts
	return p
}

// WithStartDate sets the start date for the task.
func (p *TaskParams) WithStartDate(date time.Time) *TaskParams {
	ts := date.Unix()
	p.startDate = &ts
	return p
}

// InSwimlane sets the swimlane ID for the task.
func (p *TaskParams) InSwimlane(swimlaneID int) *TaskParams {
	p.swimlaneID = &swimlaneID
	return p
}

// WithReference sets the external reference for the task.
func (p *TaskParams) WithReference(ref string) *TaskParams {
	p.reference = &ref
	return p
}

// WithTags sets the tags for the task.
func (p *TaskParams) WithTags(tags ...string) *TaskParams {
	p.tags = tags
	return p
}

// toCreateTaskRequest converts TaskParams to a CreateTaskRequest.
// The projectID is required and must be provided by the caller.
func (p *TaskParams) toCreateTaskRequest(projectID int) CreateTaskRequest {
	req := CreateTaskRequest{
		Title:     p.title,
		ProjectID: projectID,
	}

	if p.description != nil {
		req.Description = *p.description
	}
	if p.columnID != nil {
		req.ColumnID = *p.columnID
	}
	if p.categoryID != nil {
		req.CategoryID = *p.categoryID
	}
	if p.ownerID != nil {
		req.OwnerID = *p.ownerID
	}
	if p.creatorID != nil {
		req.CreatorID = *p.creatorID
	}
	if p.colorID != nil {
		req.ColorID = *p.colorID
	}
	if p.priority != nil {
		req.Priority = *p.priority
	}
	if p.score != nil {
		req.Score = *p.score
	}
	if p.dueDate != nil {
		req.DateDue = *p.dueDate
	}
	if p.startDate != nil {
		req.DateStarted = *p.startDate
	}
	if p.swimlaneID != nil {
		req.SwimlaneID = *p.swimlaneID
	}
	if p.reference != nil {
		req.Reference = *p.reference
	}
	if len(p.tags) > 0 {
		req.Tags = p.tags
	}

	return req
}
