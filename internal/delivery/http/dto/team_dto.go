package dto

import (
	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// CreateTeamRequest represents create team request body
type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Logo        string `json:"logo" binding:"omitempty,url,max=500"`
	FoundedYear int    `json:"founded_year" binding:"required,min=1800,max=2100"`
	Address     string `json:"address" binding:"omitempty,max=500"`
	City        string `json:"city" binding:"required,min=2,max=100"`
}

// UpdateTeamRequest represents update team request body
type UpdateTeamRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=255"`
	Logo        string `json:"logo" binding:"omitempty,url,max=500"`
	FoundedYear int    `json:"founded_year" binding:"omitempty,min=1800,max=2100"`
	Address     string `json:"address" binding:"omitempty,max=500"`
	City        string `json:"city" binding:"omitempty,min=2,max=100"`
}

// TeamResponse represents team data in response
type TeamResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Logo        string           `json:"logo"`
	FoundedYear int              `json:"founded_year"`
	Address     string           `json:"address"`
	City        string           `json:"city"`
	Players     []PlayerResponse `json:"players,omitempty"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

// ToTeamEntity converts CreateTeamRequest to entity.Team
func (r *CreateTeamRequest) ToTeamEntity() *entity.Team {
	return &entity.Team{
		Name:        r.Name,
		Logo:        r.Logo,
		FoundedYear: r.FoundedYear,
		Address:     r.Address,
		City:        r.City,
	}
}

// UpdateTeamEntity updates entity.Team with UpdateTeamRequest values
func (r *UpdateTeamRequest) UpdateTeamEntity(team *entity.Team) {
	if r.Name != "" {
		team.Name = r.Name
	}
	if r.Logo != "" {
		team.Logo = r.Logo
	}
	if r.FoundedYear != 0 {
		team.FoundedYear = r.FoundedYear
	}
	if r.Address != "" {
		team.Address = r.Address
	}
	if r.City != "" {
		team.City = r.City
	}
}

// ToTeamResponse converts entity.Team to TeamResponse
func ToTeamResponse(team *entity.Team) TeamResponse {
	response := TeamResponse{
		ID:          team.ID.String(),
		Name:        team.Name,
		Logo:        team.Logo,
		FoundedYear: team.FoundedYear,
		Address:     team.Address,
		City:        team.City,
		CreatedAt:   team.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   team.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if team.Players != nil {
		response.Players = make([]PlayerResponse, len(team.Players))
		for i, player := range team.Players {
			response.Players[i] = ToPlayerResponse(&player)
		}
	}

	return response
}

// ToTeamResponseList converts a slice of entity.Team to TeamResponse slice
func ToTeamResponseList(teams []entity.Team) []TeamResponse {
	responses := make([]TeamResponse, len(teams))
	for i, team := range teams {
		responses[i] = ToTeamResponse(&team)
	}
	return responses
}

// TeamSimpleResponse represents simplified team data
type TeamSimpleResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Logo string    `json:"logo"`
	City string    `json:"city"`
}

// ToTeamSimpleResponse converts entity.Team to TeamSimpleResponse
func ToTeamSimpleResponse(team *entity.Team) TeamSimpleResponse {
	return TeamSimpleResponse{
		ID:   team.ID,
		Name: team.Name,
		Logo: team.Logo,
		City: team.City,
	}
}
