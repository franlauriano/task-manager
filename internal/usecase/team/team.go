package team

import (
	"context"
	"strings"

	teamEntity "taskmanager/internal/entity/team"
	errors "taskmanager/internal/platform/errors"
	errs "taskmanager/internal/platform/errors"
	taskRepo "taskmanager/internal/repository/task"
	teamRepo "taskmanager/internal/repository/team"

	"github.com/google/uuid"
)

// Create creates a new team with business rules
func Create(ctx context.Context, t *teamEntity.Team) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.Name = strings.TrimSpace(t.Name)
	t.Description = strings.TrimSpace(t.Description)

	return teamRepo.Persist().Create(ctx, t)
}

// RetrieveByUUIDWithTasks retrieves a team by UUID with its associated tasks
func RetrieveByUUIDWithTasks(ctx context.Context, teamUUID uuid.UUID) (*teamEntity.Team, error) {
	t, err := teamRepo.Persist().RetrieveByUUID(ctx, teamUUID)
	if err != nil {
		return nil, err
	}

	tasks, err := taskRepo.Persist().ListByTeamID(ctx, t.ID)
	if err != nil {
		return nil, err
	}

	t.Tasks = tasks

	return t, nil
}

// ListPaginated lists teams with pagination
func ListPaginated(ctx context.Context, page, limit int) (*teamEntity.ListTeams, error) {
	if limit <= 0 {
		limit = Config.ListDefaultLimit
	}

	if limit > Config.ListMaxLimit {
		limit = Config.ListMaxLimit
	}

	return teamRepo.Persist().ListPaginated(ctx, page, limit)
}

// AssociateTask associates a task with a team
func AssociateTask(ctx context.Context, teamUUID, taskUUID uuid.UUID) error {
	team, err := validateAssociateTask(ctx, teamUUID, taskUUID)
	if err != nil {
		return err
	}

	return teamRepo.Persist().UpdateTaskTeamID(ctx, taskUUID, &team.ID)
}

// DisassociateTask disassociates a task from a team
func DisassociateTask(ctx context.Context, teamUUID, taskUUID uuid.UUID) error {
	err := validateDisassociateTask(ctx, teamUUID, taskUUID)
	if err != nil {
		return err
	}

	return teamRepo.Persist().UpdateTaskTeamID(ctx, taskUUID, nil)
}

// validateAssociateTask validates team and task exist and that task is not already associated with another team
func validateAssociateTask(ctx context.Context, teamUUID, taskUUID uuid.UUID) (*teamEntity.Team, error) {
	team, taskTeamID, err := validateTeamAndTask(ctx, teamUUID, taskUUID)
	if err != nil {
		return nil, err
	}

	if taskTeamID != nil && *taskTeamID != team.ID {
		return nil, &errors.ValidationErrors{Errors: []errors.ValidationError{
			{Field: "task", Message: "task is already associated with another team"},
		}}
	}

	return team, nil
}

// validateDisassociateTask validates team and task exist and that task is associated with this team
func validateDisassociateTask(ctx context.Context, teamUUID, taskUUID uuid.UUID) error {
	team, taskTeamID, err := validateTeamAndTask(ctx, teamUUID, taskUUID)
	if err != nil {
		return err
	}

	if taskTeamID == nil || *taskTeamID != team.ID {
		return &errors.ValidationErrors{Errors: []errors.ValidationError{
			{Field: "task", Message: "task is not associated with this team"},
		}}
	}

	return nil
}

// validateTeamAndTask validates that both team and task exist, returning them if valid
func validateTeamAndTask(ctx context.Context, teamUUID, taskUUID uuid.UUID) (*teamEntity.Team, *uint, error) {
	team, err := teamRepo.Persist().RetrieveByUUID(ctx, teamUUID)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, nil, &errors.ValidationErrors{Errors: []errors.ValidationError{
				{Field: "team", Message: "team not found"},
			}}
		}
		return nil, nil, err
	}

	taskTeamID, err := teamRepo.Persist().RetrieveTaskTeamID(ctx, taskUUID)
	if err != nil {
		if err == errs.ErrNotFound {
			return nil, nil, &errors.ValidationErrors{Errors: []errors.ValidationError{
				{Field: "task", Message: "task not found"},
			}}
		}
		return nil, nil, err
	}

	return team, taskTeamID, nil
}
