package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/entity"
	"github.com/zenkriztao/ayo-football-backend/internal/domain/repository"
	"gorm.io/gorm"
)

// MatchReport represents a detailed match report
type MatchReport struct {
	Match               *entity.Match             `json:"match"`
	HomeTeam            *entity.Team              `json:"home_team"`
	AwayTeam            *entity.Team              `json:"away_team"`
	HomeScore           int                       `json:"home_score"`
	AwayScore           int                       `json:"away_score"`
	MatchResult         string                    `json:"match_result"`
	MatchResultDisplay  string                    `json:"match_result_display"`
	Goals               []entity.Goal             `json:"goals"`
	TopScorer           *repository.TopScorerResult `json:"top_scorer,omitempty"`
	HomeTeamTotalWins   int64                     `json:"home_team_total_wins"`
	AwayTeamTotalWins   int64                     `json:"away_team_total_wins"`
}

// LeaderboardEntry represents a team's standings
type LeaderboardEntry struct {
	Team        *entity.Team `json:"team"`
	Played      int64        `json:"played"`
	Won         int64        `json:"won"`
	Drawn       int64        `json:"drawn"`
	Lost        int64        `json:"lost"`
	GoalsFor    int64        `json:"goals_for"`
	GoalsAgainst int64       `json:"goals_against"`
	GoalDiff    int64        `json:"goal_difference"`
	Points      int64        `json:"points"`
}

// ReportUseCase defines the interface for report operations
type ReportUseCase interface {
	GetMatchReport(ctx context.Context, matchID uuid.UUID) (*MatchReport, error)
	GetAllMatchReports(ctx context.Context, page, limit int) ([]MatchReport, int64, error)
	GetTopScorers(ctx context.Context, limit int) ([]repository.TopScorerResult, error)
}

type reportUseCaseImpl struct {
	matchRepo repository.MatchRepository
	goalRepo  repository.GoalRepository
	teamRepo  repository.TeamRepository
}

// NewReportUseCase creates a new instance of ReportUseCase
func NewReportUseCase(
	matchRepo repository.MatchRepository,
	goalRepo repository.GoalRepository,
	teamRepo repository.TeamRepository,
) ReportUseCase {
	return &reportUseCaseImpl{
		matchRepo: matchRepo,
		goalRepo:  goalRepo,
		teamRepo:  teamRepo,
	}
}

func (uc *reportUseCaseImpl) GetMatchReport(ctx context.Context, matchID uuid.UUID) (*MatchReport, error) {
	// Get match with details
	match, err := uc.matchRepo.FindByIDWithDetails(ctx, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMatchNotFound
		}
		return nil, err
	}

	// Get home team total wins
	homeWins, err := uc.matchRepo.GetTeamWinCount(ctx, match.HomeTeamID, true)
	if err != nil {
		return nil, err
	}

	// Also count away wins for home team
	homeAwayWins, err := uc.matchRepo.GetTeamWinCount(ctx, match.HomeTeamID, false)
	if err != nil {
		return nil, err
	}

	// Get away team total wins
	awayWins, err := uc.matchRepo.GetTeamWinCount(ctx, match.AwayTeamID, false)
	if err != nil {
		return nil, err
	}

	// Also count home wins for away team
	awayHomeWins, err := uc.matchRepo.GetTeamWinCount(ctx, match.AwayTeamID, true)
	if err != nil {
		return nil, err
	}

	// Get top scorer
	topScorers, err := uc.goalRepo.GetTopScorers(ctx, 1)
	if err != nil {
		return nil, err
	}

	var topScorer *repository.TopScorerResult
	if len(topScorers) > 0 {
		topScorer = &topScorers[0]
	}

	homeScore := 0
	awayScore := 0
	if match.HomeScore != nil {
		homeScore = *match.HomeScore
	}
	if match.AwayScore != nil {
		awayScore = *match.AwayScore
	}

	report := &MatchReport{
		Match:              match,
		HomeTeam:           match.HomeTeam,
		AwayTeam:           match.AwayTeam,
		HomeScore:          homeScore,
		AwayScore:          awayScore,
		MatchResult:        string(match.GetResult()),
		MatchResultDisplay: match.GetResultDisplay(),
		Goals:              match.Goals,
		TopScorer:          topScorer,
		HomeTeamTotalWins:  homeWins + homeAwayWins,
		AwayTeamTotalWins:  awayWins + awayHomeWins,
	}

	return report, nil
}

func (uc *reportUseCaseImpl) GetAllMatchReports(ctx context.Context, page, limit int) ([]MatchReport, int64, error) {
	matches, total, err := uc.matchRepo.GetCompletedMatches(ctx, page, limit)
	if err != nil {
		return nil, 0, err
	}

	reports := make([]MatchReport, len(matches))
	for i, match := range matches {
		homeScore := 0
		awayScore := 0
		if match.HomeScore != nil {
			homeScore = *match.HomeScore
		}
		if match.AwayScore != nil {
			awayScore = *match.AwayScore
		}

		// Get team win counts
		homeWins, _ := uc.matchRepo.GetTeamWinCount(ctx, match.HomeTeamID, true)
		homeAwayWins, _ := uc.matchRepo.GetTeamWinCount(ctx, match.HomeTeamID, false)
		awayWins, _ := uc.matchRepo.GetTeamWinCount(ctx, match.AwayTeamID, false)
		awayHomeWins, _ := uc.matchRepo.GetTeamWinCount(ctx, match.AwayTeamID, true)

		reports[i] = MatchReport{
			Match:              &matches[i],
			HomeTeam:           match.HomeTeam,
			AwayTeam:           match.AwayTeam,
			HomeScore:          homeScore,
			AwayScore:          awayScore,
			MatchResult:        string(match.GetResult()),
			MatchResultDisplay: match.GetResultDisplay(),
			Goals:              match.Goals,
			HomeTeamTotalWins:  homeWins + homeAwayWins,
			AwayTeamTotalWins:  awayWins + awayHomeWins,
		}
	}

	// Get top scorer for overall
	topScorers, err := uc.goalRepo.GetTopScorers(ctx, 1)
	if err == nil && len(topScorers) > 0 && len(reports) > 0 {
		reports[0].TopScorer = &topScorers[0]
	}

	return reports, total, nil
}

func (uc *reportUseCaseImpl) GetTopScorers(ctx context.Context, limit int) ([]repository.TopScorerResult, error) {
	return uc.goalRepo.GetTopScorers(ctx, limit)
}
