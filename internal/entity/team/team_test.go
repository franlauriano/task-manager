package team

import (
	"strings"
	"testing"

	errors "taskmanager/internal/platform/errors"
	"taskmanager/internal/platform/testing/assert"
)

func TestTeam_Validate(t *testing.T) {
	tests := []struct {
		name    string
		team    *Team
		wantErr *errors.ValidationErrors
	}{
		{
			"Validate team with success",
			&Team{
				Name:        "Development Team",
				Description: "Team responsible for development",
			},
			nil,
		},
		{
			"Validate team with empty name",
			&Team{
				Name:        "",
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
				},
			},
		},
		{
			"Validate team with empty description",
			&Team{
				Name:        "Valid name",
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
			"Validate team with only whitespace name",
			&Team{
				Name:        "   ",
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
				},
			},
		},
		{
			"Validate team with only whitespace description",
			&Team{
				Name:        "Valid name",
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
			"Validate team with empty name and description",
			&Team{
				Name:        "",
				Description: "",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate team with only whitespace name and description",
			&Team{
				Name:        "   ",
				Description: "   ",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
					{
						Field:   "description",
						Message: "description is required",
					},
				},
			},
		},
		{
			"Validate team with tab and newline whitespace in name",
			&Team{
				Name:        "\t\n\r   ",
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "name",
						Message: "name is required",
					},
				},
			},
		},
		{
			"Validate team with tab and newline whitespace in description",
			&Team{
				Name:        "Valid name",
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
			"Validate team with name exceeding 255 characters",
			&Team{
				Name:        strings.Repeat("a", 256),
				Description: "Valid description",
			},
			&errors.ValidationErrors{
				Errors: []errors.ValidationError{
					{
						Field:   "name",
						Message: "name must not exceed 255 characters",
					},
				},
			},
		},
		{
			"Validate team with name exactly 255 characters",
			&Team{
				Name:        strings.Repeat("a", 255),
				Description: "Valid description",
			},
			nil,
		},
		{
			"Validate team with name with leading and trailing whitespace",
			&Team{
				Name:        "  Valid name  ",
				Description: "Valid description",
			},
			nil,
		},
		{
			"Validate team with description with leading and trailing whitespace",
			&Team{
				Name:        "Valid name",
				Description: "  Valid description  ",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.team.Validate()
			if diff := assert.CompareErrors(err, tt.wantErr); diff != "" {
				t.Errorf("Team.Validate() error diff: %s", diff)
				return
			}
		})
	}
}
