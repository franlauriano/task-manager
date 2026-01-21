package dto

import "taskmanager/internal/entity/team"

// CreateTeamRequest represents the payload for creating a new team
type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToTeam converts CreateTeamRequest to team.Team
func (r *CreateTeamRequest) ToTeam() *team.Team {
	return &team.Team{
		Name:        r.Name,
		Description: r.Description,
	}
}
