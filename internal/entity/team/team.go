package team

import (
	"strings"

	taskEntity "taskmanager/internal/entity/task"
	errors "taskmanager/internal/platform/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Team represents a team entity
type Team struct {
	gorm.Model

	UUID        uuid.UUID         `gorm:"type:uuid;uniqueIndex;not null" json:"-"`
	Name        string            `gorm:"not null" json:"-"`
	Description string            `gorm:"not null" json:"-"`
	Tasks       []taskEntity.Task `gorm:"foreignKey:TeamID;references:ID" json:"-"`
}

// ListTeams contains paginated teams and total count
type ListTeams struct {
	Teams      []Team
	TotalItems int
	Limit      int
	Page       int
}

// BeforeCreate is a GORM hook to generate UUID v7 before creating
func (t *Team) BeforeCreate(tx *gorm.DB) (err error) {
	if t.UUID == (uuid.UUID{}) {
		t.UUID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return nil
}

// AfterFind is a GORM hook to normalize timestamps
func (t *Team) AfterFind(tx *gorm.DB) (err error) {
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

// Validate validates the team fields
func (t *Team) Validate() *errors.ValidationErrors {
	var errs []errors.ValidationError

	name := strings.TrimSpace(t.Name)
	if name == "" {
		errs = append(errs, errors.ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	} else if len(name) > 255 {
		errs = append(errs, errors.ValidationError{
			Field:   "name",
			Message: "name must not exceed 255 characters",
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
