package team

import (
	"context"
	"errors"

	"taskmanager/internal/entity/team"
	"taskmanager/internal/platform/database"
	errs "taskmanager/internal/platform/errors"
	"taskmanager/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Persistent defines the interface for team persistence
type Persistent interface {
	Create(ctx context.Context, t *team.Team) error
	RetrieveByUUID(ctx context.Context, teamUUID uuid.UUID) (*team.Team, error)
	ListPaginated(ctx context.Context, page, limit int) (*team.ListTeams, error)
	RetrieveTaskTeamID(ctx context.Context, taskUUID uuid.UUID) (*uint, error)
	UpdateTaskTeamID(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error
}

// datasource implements the persistent interface using PostgreSQL
type datasource struct {
	alias string
}

var persist Persistent = &datasource{alias: repository.DatabaseAlias}

// SetPersist sets the persistent implementation
func SetPersist(p Persistent) {
	persist = p
}

// Persist returns the current persistent implementation
func Persist() Persistent {
	return persist
}

// Create saves a new team to the database
func (p *datasource) Create(ctx context.Context, t *team.Team) error {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return err
	}

	if err := db.Create(t).Error; err != nil {
		return err
	}

	return nil
}

// RetrieveByUUID retrieves a team by UUID from the database
func (p *datasource) RetrieveByUUID(ctx context.Context, teamUUID uuid.UUID) (*team.Team, error) {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return nil, err
	}

	var t team.Team
	if err := db.Where("uuid = ?", teamUUID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

// ListPaginated lists teams with pagination from the database
func (p *datasource) ListPaginated(ctx context.Context, page, limit int) (*team.ListTeams, error) {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return nil, err
	}

	var teams []team.Team
	var totalItems int64

	query := db.Model(&team.Team{})

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Order("id DESC").Offset(offset).Limit(limit).Find(&teams).Error; err != nil {
		return nil, err
	}

	return &team.ListTeams{
		Limit:      limit,
		Page:       page,
		Teams:      teams,
		TotalItems: int(totalItems),
	}, nil
}

// RetrieveTaskTeamID retrieves the team_id of a task by UUID
func (p *datasource) RetrieveTaskTeamID(ctx context.Context, taskUUID uuid.UUID) (*uint, error) {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return nil, err
	}

	var result struct {
		TeamID *uint `gorm:"column:team_id"`
	}
	if err := db.Table("tasks").
		Select("team_id").
		Where("uuid = ? AND deleted_at IS NULL", taskUUID).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return result.TeamID, nil
}

// UpdateTaskTeamID updates the team_id of a task
func (p *datasource) UpdateTaskTeamID(ctx context.Context, taskUUID uuid.UUID, teamID *uint) error {
	db, err := database.DBFromContext(ctx, p.alias)
	if err != nil {
		return err
	}

	result := db.Table("tasks").
		Where("uuid = ?", taskUUID).
		Update("team_id", teamID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
