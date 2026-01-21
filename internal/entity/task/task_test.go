package task

import (
	"strings"
	"testing"
	"time"

	errors "taskmanager/internal/platform/errors"
	"taskmanager/internal/platform/testing/assert"

	"github.com/google/go-cmp/cmp"
)

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    *Task
		wantErr *errors.ValidationErrors
	}{
		{
			"Validate task with success",
			&Task{
				Title:       "Complete project documentation",
				Description: "Write comprehensive documentation for the project",
			},
			nil,
		},
		{
			"Validate task with empty title",
			&Task{
				Title:       "",
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Validate task with empty description",
			&Task{
				Title:       "Valid title",
				Description: "",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate task with only whitespace title",
			&Task{
				Title:       "   ",
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Validate task with only whitespace description",
			&Task{
				Title:       "Valid title",
				Description: "   ",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate task with empty title and description",
			&Task{
				Title:       "",
				Description: "",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate task with only whitespace title and description",
			&Task{
				Title:       "   ",
				Description: "   ",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate task with tab and newline whitespace in title",
			&Task{
				Title:       "\t\n\r   ",
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "title",
						Message: "title is required",
					},
				},
			},
		},
		{
			"Validate task with tab and newline whitespace in description",
			&Task{
				Title:       "Valid title",
				Description: "\t\n\r   ",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate task with title exceeding 255 characters",
			&Task{
				Title:       strings.Repeat("a", 256),
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "title",
						Message: "title must not exceed 255 characters",
					},
				},
			},
		},
		{
			"Validate task with title exactly 255 characters",
			&Task{
				Title:       strings.Repeat("a", 255),
				Description: "Valid description",
			},
			nil,
		},
		{
			"Validate task with title with leading and trailing whitespace",
			&Task{
				Title:       "  Valid title  ",
				Description: "Valid description",
			},
			nil,
		},
		{
			"Validate task with description with leading and trailing whitespace",
			&Task{
				Title:       "Valid title",
				Description: "  Valid description  ",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("Task.Validate() error diff: %s", diff)
				return
			}
		})
	}
}

func TestTaskStatus_ValidateTransitionTo(t *testing.T) {
	tests := []struct {
		name      string
		status    TaskStatus
		newStatus TaskStatus
		wantErr   error
	}{
		{
			"Validate transition from todo to in_progress",
			StatusTodo,
			StatusInProgress,
			nil,
		},
		{
			"Validate transition from todo to canceled",
			StatusTodo,
			StatusCanceled,
			nil,
		},
		{
			"Validate transition from in_progress to canceled",
			StatusInProgress,
			StatusCanceled,
			nil,
		},
		{
			"Validate transition from in_progress to done",
			StatusInProgress,
			StatusDone,
			nil,
		},
		{
			"Validate transition from todo to done - invalid",
			StatusTodo,
			StatusDone,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from todo to todo - invalid",
			StatusTodo,
			StatusTodo,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from in_progress to todo - invalid",
			StatusInProgress,
			StatusTodo,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from in_progress to in_progress - invalid",
			StatusInProgress,
			StatusInProgress,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from canceled to any status - invalid",
			StatusCanceled,
			StatusTodo,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from canceled to in_progress - invalid",
			StatusCanceled,
			StatusInProgress,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from canceled to done - invalid",
			StatusCanceled,
			StatusDone,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from done to any status - invalid",
			StatusDone,
			StatusTodo,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from done to in_progress - invalid",
			StatusDone,
			StatusInProgress,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition from done to canceled - invalid",
			StatusDone,
			StatusCanceled,
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status transition",
					},
				},
			},
		},
		{
			"Validate transition to invalid status",
			StatusTodo,
			TaskStatus("invalid_status"),
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status value",
					},
				},
			},
		},
		{
			"Validate transition from in_progress to invalid status",
			StatusInProgress,
			TaskStatus("invalid_status"),
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "status",
						Message: "invalid status value",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.status.ValidateTransitionTo(tt.newStatus)
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("TaskStatus.ValidateTransitionTo() error diff: %s", diff)
				return
			}
		})
	}
}

func TestTask_EnsureTimestampsForStatus(t *testing.T) {
	existingStartedTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	existingFinishedTime := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)
	testTimestamp := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		task      *Task
		newStatus TaskStatus
		timestamp *time.Time
		want      func() (*time.Time, *time.Time)
	}{
		{
			"Ensure StartedAt is set when transitioning to in_progress with nil StartedAt",
			&Task{
				StartedAt:  nil,
				FinishedAt: nil,
			},
			StatusInProgress,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				// StartedAt should be set to the timestamp, FinishedAt should remain nil
				return &testTimestamp, nil
			},
		},
		{
			"Ensure StartedAt is preserved when transitioning to in_progress with existing StartedAt",
			&Task{
				StartedAt:  &existingStartedTime,
				FinishedAt: nil,
			},
			StatusInProgress,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				return &existingStartedTime, nil
			},
		},
		{
			"Ensure FinishedAt is set when transitioning to done with nil FinishedAt",
			&Task{
				StartedAt:  nil,
				FinishedAt: nil,
			},
			StatusDone,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				// StartedAt should remain nil, FinishedAt should be set to the timestamp
				return nil, &testTimestamp
			},
		},
		{
			"Ensure FinishedAt is preserved when transitioning to done with existing FinishedAt",
			&Task{
				StartedAt:  nil,
				FinishedAt: &existingFinishedTime,
			},
			StatusDone,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				return nil, &existingFinishedTime
			},
		},
		{
			"Ensure FinishedAt is set when transitioning to canceled with nil FinishedAt",
			&Task{
				StartedAt:  nil,
				FinishedAt: nil,
			},
			StatusCanceled,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				// StartedAt should remain nil, FinishedAt should be set to the timestamp
				return nil, &testTimestamp
			},
		},
		{
			"Ensure FinishedAt is preserved when transitioning to canceled with existing FinishedAt",
			&Task{
				StartedAt:  nil,
				FinishedAt: &existingFinishedTime,
			},
			StatusCanceled,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				return nil, &existingFinishedTime
			},
		},
		{
			"Ensure no timestamps are set when transitioning to todo",
			&Task{
				StartedAt:  nil,
				FinishedAt: nil,
			},
			StatusTodo,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				return nil, nil
			},
		},
		{
			"Ensure timestamps are not modified when transitioning to todo with existing timestamps",
			&Task{
				StartedAt:  &existingStartedTime,
				FinishedAt: &existingFinishedTime,
			},
			StatusTodo,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				return &existingStartedTime, &existingFinishedTime
			},
		},
		{
			"Ensure FinishedAt is set when transitioning to canceled with existing StartedAt",
			&Task{
				StartedAt:  &existingStartedTime,
				FinishedAt: nil,
			},
			StatusCanceled,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				// StartedAt should be preserved, FinishedAt should be set to the timestamp
				return &existingStartedTime, &testTimestamp
			},
		},
		{
			"Ensure StartedAt is not set when transitioning to done",
			&Task{
				StartedAt:  nil,
				FinishedAt: nil,
			},
			StatusDone,
			&testTimestamp,
			func() (*time.Time, *time.Time) {
				// StartedAt should remain nil, FinishedAt should be set to the timestamp
				return nil, &testTimestamp
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			wantStartedAt, wantFinishedAt := tt.want()
			tt.task.EnsureTimestampsForStatus(tt.newStatus, tt.timestamp)

			if diff := cmp.Diff(tt.task.StartedAt, wantStartedAt); diff != "" {
				t.Errorf("Task.EnsureTimestampsForStatus() StartedAt diff: %s", diff)
				return
			}

			if diff := cmp.Diff(tt.task.FinishedAt, wantFinishedAt); diff != "" {
				t.Errorf("Task.EnsureTimestampsForStatus() FinishedAt diff: %s", diff)
			}
		})
	}
}
