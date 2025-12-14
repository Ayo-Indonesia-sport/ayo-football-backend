package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
)

// CreateMatchRequest represents create match request body
type CreateMatchRequest struct {
	MatchDate  string `json:"match_date" binding:"required"` // Format: 2006-01-02
	MatchTime  string `json:"match_time" binding:"required"` // Format: 15:04
	HomeTeamID string `json:"home_team_id" binding:"required,uuid"`
	AwayTeamID string `json:"away_team_id" binding:"required,uuid"`
}

// UpdateMatchRequest represents update match request body
type UpdateMatchRequest struct {
	MatchDate  string `json:"match_date" binding:"omitempty"` // Format: 2006-01-02
	MatchTime  string `json:"match_time" binding:"omitempty"` // Format: 15:04
	HomeTeamID string `json:"home_team_id" binding:"omitempty,uuid"`
	AwayTeamID string `json:"away_team_id" binding:"omitempty,uuid"`
	Status     string `json:"status" binding:"omitempty,oneof=scheduled ongoing completed cancelled"`
}

// RecordMatchResultRequest represents match result recording request body
type RecordMatchResultRequest struct {
	HomeScore int        `json:"home_score" binding:"min=0"`
	AwayScore int        `json:"away_score" binding:"min=0"`
	Goals     []GoalRequest `json:"goals" binding:"dive"`
}

// GoalRequest represents a goal input
type GoalRequest struct {
	PlayerID  string `json:"player_id" binding:"required,uuid"`
	TeamID    string `json:"team_id" binding:"required,uuid"`
	Minute    int    `json:"minute" binding:"required,min=1,max=120"`
	IsOwnGoal bool   `json:"is_own_goal"`
}

// MatchResponse represents match data in response
type MatchResponse struct {
	ID           string              `json:"id"`
	MatchDate    string              `json:"match_date"`
	MatchTime    string              `json:"match_time"`
	HomeTeamID   string              `json:"home_team_id"`
	AwayTeamID   string              `json:"away_team_id"`
	HomeScore    *int                `json:"home_score"`
	AwayScore    *int                `json:"away_score"`
	Status       string              `json:"status"`
	StatusName   string              `json:"status_name"`
	HomeTeam     *TeamSimpleResponse `json:"home_team,omitempty"`
	AwayTeam     *TeamSimpleResponse `json:"away_team,omitempty"`
	Goals        []GoalResponse      `json:"goals,omitempty"`
	MatchResult  string              `json:"match_result,omitempty"`
	ResultDisplay string             `json:"result_display,omitempty"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
}

// GoalResponse represents goal data in response
type GoalResponse struct {
	ID         string `json:"id"`
	MatchID    string `json:"match_id"`
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name,omitempty"`
	TeamID     string `json:"team_id"`
	TeamName   string `json:"team_name,omitempty"`
	Minute     int    `json:"minute"`
	IsOwnGoal  bool   `json:"is_own_goal"`
}

// ToMatchEntity converts CreateMatchRequest to entity.Match
func (r *CreateMatchRequest) ToMatchEntity() (*entity.Match, error) {
	homeTeamID, err := uuid.Parse(r.HomeTeamID)
	if err != nil {
		return nil, err
	}

	awayTeamID, err := uuid.Parse(r.AwayTeamID)
	if err != nil {
		return nil, err
	}

	matchDate, err := time.Parse("2006-01-02", r.MatchDate)
	if err != nil {
		return nil, err
	}

	return &entity.Match{
		MatchDate:  matchDate,
		MatchTime:  r.MatchTime,
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
		Status:     entity.MatchStatusScheduled,
	}, nil
}

// UpdateMatchEntity updates entity.Match with UpdateMatchRequest values
func (r *UpdateMatchRequest) UpdateMatchEntity(match *entity.Match) error {
	if r.MatchDate != "" {
		matchDate, err := time.Parse("2006-01-02", r.MatchDate)
		if err != nil {
			return err
		}
		match.MatchDate = matchDate
	}
	if r.MatchTime != "" {
		match.MatchTime = r.MatchTime
	}
	if r.HomeTeamID != "" {
		homeTeamID, err := uuid.Parse(r.HomeTeamID)
		if err != nil {
			return err
		}
		match.HomeTeamID = homeTeamID
	}
	if r.AwayTeamID != "" {
		awayTeamID, err := uuid.Parse(r.AwayTeamID)
		if err != nil {
			return err
		}
		match.AwayTeamID = awayTeamID
	}
	if r.Status != "" {
		match.Status = entity.MatchStatus(r.Status)
	}
	return nil
}

// ToMatchResponse converts entity.Match to MatchResponse
func ToMatchResponse(match *entity.Match) MatchResponse {
	response := MatchResponse{
		ID:            match.ID.String(),
		MatchDate:     match.MatchDate.Format("2006-01-02"),
		MatchTime:     match.MatchTime,
		HomeTeamID:    match.HomeTeamID.String(),
		AwayTeamID:    match.AwayTeamID.String(),
		HomeScore:     match.HomeScore,
		AwayScore:     match.AwayScore,
		Status:        string(match.Status),
		StatusName:    getMatchStatusDisplayName(match.Status),
		MatchResult:   string(match.GetResult()),
		ResultDisplay: match.GetResultDisplay(),
		CreatedAt:     match.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     match.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if match.HomeTeam != nil {
		homeTeam := ToTeamSimpleResponse(match.HomeTeam)
		response.HomeTeam = &homeTeam
	}

	if match.AwayTeam != nil {
		awayTeam := ToTeamSimpleResponse(match.AwayTeam)
		response.AwayTeam = &awayTeam
	}

	if match.Goals != nil {
		response.Goals = ToGoalResponseList(match.Goals)
	}

	return response
}

// ToMatchResponseList converts a slice of entity.Match to MatchResponse slice
func ToMatchResponseList(matches []entity.Match) []MatchResponse {
	responses := make([]MatchResponse, len(matches))
	for i, match := range matches {
		responses[i] = ToMatchResponse(&match)
	}
	return responses
}

// ToGoalResponse converts entity.Goal to GoalResponse
func ToGoalResponse(goal *entity.Goal) GoalResponse {
	response := GoalResponse{
		ID:        goal.ID.String(),
		MatchID:   goal.MatchID.String(),
		PlayerID:  goal.PlayerID.String(),
		TeamID:    goal.TeamID.String(),
		Minute:    goal.Minute,
		IsOwnGoal: goal.IsOwnGoal,
	}

	if goal.Player != nil {
		response.PlayerName = goal.Player.Name
	}

	if goal.Team != nil {
		response.TeamName = goal.Team.Name
	}

	return response
}

// ToGoalResponseList converts a slice of entity.Goal to GoalResponse slice
func ToGoalResponseList(goals []entity.Goal) []GoalResponse {
	responses := make([]GoalResponse, len(goals))
	for i, goal := range goals {
		responses[i] = ToGoalResponse(&goal)
	}
	return responses
}

// getMatchStatusDisplayName returns the display name for match status
func getMatchStatusDisplayName(status entity.MatchStatus) string {
	switch status {
	case entity.MatchStatusScheduled:
		return "Scheduled"
	case entity.MatchStatusOngoing:
		return "Ongoing"
	case entity.MatchStatusCompleted:
		return "Completed"
	case entity.MatchStatusCancelled:
		return "Cancelled"
	default:
		return string(status)
	}
}
