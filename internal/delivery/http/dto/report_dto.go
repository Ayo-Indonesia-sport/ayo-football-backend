package dto

import (
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/usecase"
)

// MatchReportResponse represents match report data in response
type MatchReportResponse struct {
	Match              MatchResponse       `json:"match"`
	HomeTeam           TeamSimpleResponse  `json:"home_team"`
	AwayTeam           TeamSimpleResponse  `json:"away_team"`
	HomeScore          int                 `json:"home_score"`
	AwayScore          int                 `json:"away_score"`
	MatchResult        string              `json:"match_result"`
	MatchResultDisplay string              `json:"match_result_display"`
	Goals              []GoalResponse      `json:"goals"`
	TopScorer          *TopScorerResponse  `json:"top_scorer,omitempty"`
	HomeTeamTotalWins  int64               `json:"home_team_total_wins"`
	AwayTeamTotalWins  int64               `json:"away_team_total_wins"`
}

// TopScorerResponse represents top scorer data in response
type TopScorerResponse struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	TeamID     string `json:"team_id"`
	TeamName   string `json:"team_name"`
	GoalCount  int64  `json:"goal_count"`
}

// ToMatchReportResponse converts usecase.MatchReport to MatchReportResponse
func ToMatchReportResponse(report *usecase.MatchReport) MatchReportResponse {
	response := MatchReportResponse{
		Match:              ToMatchResponse(report.Match),
		HomeScore:          report.HomeScore,
		AwayScore:          report.AwayScore,
		MatchResult:        report.MatchResult,
		MatchResultDisplay: report.MatchResultDisplay,
		HomeTeamTotalWins:  report.HomeTeamTotalWins,
		AwayTeamTotalWins:  report.AwayTeamTotalWins,
	}

	if report.HomeTeam != nil {
		response.HomeTeam = ToTeamSimpleResponse(report.HomeTeam)
	}

	if report.AwayTeam != nil {
		response.AwayTeam = ToTeamSimpleResponse(report.AwayTeam)
	}

	if report.Goals != nil {
		response.Goals = ToGoalResponseList(report.Goals)
	}

	if report.TopScorer != nil {
		response.TopScorer = ToTopScorerResponse(report.TopScorer)
	}

	return response
}

// ToMatchReportResponseList converts a slice of usecase.MatchReport to MatchReportResponse slice
func ToMatchReportResponseList(reports []usecase.MatchReport) []MatchReportResponse {
	responses := make([]MatchReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = ToMatchReportResponse(&report)
	}
	return responses
}

// ToTopScorerResponse converts repository.TopScorerResult to TopScorerResponse
func ToTopScorerResponse(scorer *repository.TopScorerResult) *TopScorerResponse {
	if scorer == nil {
		return nil
	}
	return &TopScorerResponse{
		PlayerID:   scorer.PlayerID.String(),
		PlayerName: scorer.PlayerName,
		TeamID:     scorer.TeamID.String(),
		TeamName:   scorer.TeamName,
		GoalCount:  scorer.GoalCount,
	}
}

// ToTopScorerResponseList converts a slice of repository.TopScorerResult to TopScorerResponse slice
func ToTopScorerResponseList(scorers []repository.TopScorerResult) []TopScorerResponse {
	responses := make([]TopScorerResponse, len(scorers))
	for i, scorer := range scorers {
		responses[i] = *ToTopScorerResponse(&scorer)
	}
	return responses
}
