package task

import (
	"strings"
	"taskmanager/internal/platform/errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "to_do"
	StatusInProgress TaskStatus = "in_progress"
	StatusCanceled   TaskStatus = "canceled"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	gorm.Model

	UUID        uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"-"`
	Title       string     `gorm:"not null" json:"-"`
	Description string     `gorm:"not null" json:"-"`
	Status      TaskStatus `gorm:"type:varchar(20);not null;default:'todo'" json:"-"`
	FinishedAt  *time.Time `json:"-"`
	StartedAt   *time.Time `json:"-"`
	TeamID      *uint      `gorm:"index" json:"-"`
}

// ListTasks contains paginated tasks and total count
type ListTasks struct {
	Tasks      []Task
	TotalItems int
	Limit      int
	Page       int
}

// BeforeCreate is a GORM hook to generate UUID v7 before creating
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.UUID == (uuid.UUID{}) {
		t.UUID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return nil
}

// AfterFind is a GORM hook to normalize timestamps
func (t *Task) AfterFind(tx *gorm.DB) (err error) {
	if !t.CreatedAt.IsZero() {
		t.CreatedAt = t.CreatedAt.UTC()
	}
	if !t.UpdatedAt.IsZero() {
		t.UpdatedAt = t.UpdatedAt.UTC()
	}
	if t.DeletedAt.Valid && !t.DeletedAt.Time.IsZero() {
		t.DeletedAt.Time = t.DeletedAt.Time.UTC()
	}
	return nil
}

// Validate validates the task fields
func (t *Task) Validate() *errors.ValidationErrors {
	var errs []errors.ValidationError

	title := strings.TrimSpace(t.Title)
	if title == "" {
		errs = append(errs, errors.ValidationError{
			Field:   "title",
			Message: "title is required",
		})
	} else if len(title) > 255 {
		errs = append(errs, errors.ValidationError{
			Field:   "title",
			Message: "title must not exceed 255 characters",
		})
	}

	description := strings.TrimSpace(t.Description)
	if description == "" {
		errs = append(errs, errors.ValidationError{
			Field:   "description",
			Message: "description is required",
		})
	}

	if len(errs) > 0 {
		return &errors.ValidationErrors{Errors: errs}
	}

	return nil
}

// ValidateTransitionTo validates if the status transition is allowed
func (status TaskStatus) ValidateTransitionTo(new TaskStatus) error {
	if !new.validate() {
		return &errors.ValidationErrors{Errors: []errors.ValidationError{
			{Field: "status", Message: "invalid status value"},
		}}
	}

	allowedTransitions := map[TaskStatus][]TaskStatus{
		StatusTodo:       {StatusInProgress, StatusCanceled},
		StatusInProgress: {StatusCanceled, StatusDone},
	}

	allowed, exists := allowedTransitions[status]
	if !exists {
		return &errors.ValidationErrors{Errors: []errors.ValidationError{
			{Field: "status", Message: "invalid status transition"},
		}}
	}

	for _, status := range allowed {
		if status == new {
			return nil
		}
	}

	return &errors.ValidationErrors{Errors: []errors.ValidationError{
		{Field: "status", Message: "invalid status transition"},
	}}
}

// EnsureTimestampsForStatus ensures timestamps are set based on the new status
// It only sets timestamps if they are nil, preserving existing values
func (t *Task) EnsureTimestampsForStatus(newStatus TaskStatus, timestamp *time.Time) {
	switch newStatus {
	case StatusInProgress:
		if t.StartedAt == nil {
			t.StartedAt = timestamp
		}
	case StatusDone, StatusCanceled:
		if t.FinishedAt == nil {
			t.FinishedAt = timestamp
		}
	}
}

// validate validates if the status is valid
func (status TaskStatus) validate() bool {
	validStatuses := []TaskStatus{
		StatusTodo,
		StatusInProgress,
		StatusCanceled,
		StatusDone,
	}

	for _, vs := range validStatuses {
		if status == vs {
			return true
		}
	}

	return false
}
