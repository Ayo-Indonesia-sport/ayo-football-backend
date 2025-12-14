package dto

import (
	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// CreatePlayerRequest represents create player request body
type CreatePlayerRequest struct {
	TeamID       string `json:"team_id" binding:"required,uuid"`
	Name         string `json:"name" binding:"required,min=2,max=255"`
	Height       float64 `json:"height" binding:"required,min=100,max=250"`
	Weight       float64 `json:"weight" binding:"required,min=30,max=200"`
	Position     string `json:"position" binding:"required,oneof=forward midfielder defender goalkeeper"`
	JerseyNumber int    `json:"jersey_number" binding:"required,min=1,max=99"`
}

// UpdatePlayerRequest represents update player request body
type UpdatePlayerRequest struct {
	TeamID       string  `json:"team_id" binding:"omitempty,uuid"`
	Name         string  `json:"name" binding:"omitempty,min=2,max=255"`
	Height       float64 `json:"height" binding:"omitempty,min=100,max=250"`
	Weight       float64 `json:"weight" binding:"omitempty,min=30,max=200"`
	Position     string  `json:"position" binding:"omitempty,oneof=forward midfielder defender goalkeeper"`
	JerseyNumber int     `json:"jersey_number" binding:"omitempty,min=1,max=99"`
}

// PlayerResponse represents player data in response
type PlayerResponse struct {
	ID           string             `json:"id"`
	TeamID       string             `json:"team_id"`
	Name         string             `json:"name"`
	Height       float64            `json:"height"`
	Weight       float64            `json:"weight"`
	Position     string             `json:"position"`
	PositionName string             `json:"position_name"`
	JerseyNumber int                `json:"jersey_number"`
	Team         *TeamSimpleResponse `json:"team,omitempty"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}

// ToPlayerEntity converts CreatePlayerRequest to entity.Player
func (r *CreatePlayerRequest) ToPlayerEntity() (*entity.Player, error) {
	teamID, err := uuid.Parse(r.TeamID)
	if err != nil {
		return nil, err
	}

	return &entity.Player{
		TeamID:       teamID,
		Name:         r.Name,
		Height:       r.Height,
		Weight:       r.Weight,
		Position:     entity.PlayerPosition(r.Position),
		JerseyNumber: r.JerseyNumber,
	}, nil
}

// UpdatePlayerEntity updates entity.Player with UpdatePlayerRequest values
func (r *UpdatePlayerRequest) UpdatePlayerEntity(player *entity.Player) error {
	if r.TeamID != "" {
		teamID, err := uuid.Parse(r.TeamID)
		if err != nil {
			return err
		}
		player.TeamID = teamID
	}
	if r.Name != "" {
		player.Name = r.Name
	}
	if r.Height != 0 {
		player.Height = r.Height
	}
	if r.Weight != 0 {
		player.Weight = r.Weight
	}
	if r.Position != "" {
		player.Position = entity.PlayerPosition(r.Position)
	}
	if r.JerseyNumber != 0 {
		player.JerseyNumber = r.JerseyNumber
	}
	return nil
}

// ToPlayerResponse converts entity.Player to PlayerResponse
func ToPlayerResponse(player *entity.Player) PlayerResponse {
	response := PlayerResponse{
		ID:           player.ID.String(),
		TeamID:       player.TeamID.String(),
		Name:         player.Name,
		Height:       player.Height,
		Weight:       player.Weight,
		Position:     string(player.Position),
		PositionName: getPositionDisplayName(player.Position),
		JerseyNumber: player.JerseyNumber,
		CreatedAt:    player.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    player.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if player.Team != nil {
		teamSimple := ToTeamSimpleResponse(player.Team)
		response.Team = &teamSimple
	}

	return response
}

// ToPlayerResponseList converts a slice of entity.Player to PlayerResponse slice
func ToPlayerResponseList(players []entity.Player) []PlayerResponse {
	responses := make([]PlayerResponse, len(players))
	for i, player := range players {
		responses[i] = ToPlayerResponse(&player)
	}
	return responses
}

// getPositionDisplayName returns the display name for a position
func getPositionDisplayName(position entity.PlayerPosition) string {
	switch position {
	case entity.PositionForward:
		return "Penyerang"
	case entity.PositionMidfielder:
		return "Gelandang"
	case entity.PositionDefender:
		return "Bertahan"
	case entity.PositionGoalkeeper:
		return "Penjaga Gawang"
	default:
		return string(position)
	}
}
